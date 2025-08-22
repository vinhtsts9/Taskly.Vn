package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"Taskly.com/m/global"
	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type paymentService struct {
	store database.Store // Sử dụng store interface từ sqlc
}

// NewPaymentService khởi tạo một payment service mới.
func NewPaymentService(store database.Store) service.IPaymentService {
	return &paymentService{
		store: store,
	}
}

func (s *paymentService) CreatePaymentIntent(
    ctx context.Context,
    params model.CreatePaymentParams,
    userInfo uuid.UUID,
    ipAddr string,
    idempotencyKey string, // thêm param từ header
) (*model.VNPayPaymentURLResponse, error) {

    // 0. Validate basic
    if params.OrderID == uuid.Nil {
        return nil, fmt.Errorf("order id required")
    }

    // 1. Get order
    order, err := s.store.GetOrderByID(ctx, params.OrderID)
    if err != nil {
        return nil, errors.New("order not found")
    }
    if order.BuyerID != userInfo {
        return nil, errors.New("you are not authorized to pay for this order")
    }
    if order.TotalPrice <= 0 {
        return nil, errors.New("invalid order amount")
    }
    if order.Status != "pending" && order.Status != "active" {
        return nil, fmt.Errorf("order cannot be paid in status %s", order.Status)
    }

    // 2. Idempotency: nếu có idempotencyKey -> tìm topup đã tạo trước đó (chống double-click) và tìm theo order để tránh tạo mới topup
	topup, err := s.store.GetTopupByIdempotencyKeyAndOrder(ctx, database.GetTopupByIdempotencyKeyAndOrderParams{
		IdempotencyKey: idempotencyKey,
		OrderID:        utils.WrapUUID(params.OrderID),
	})
	if err == nil {
		if topup.Status == "PENDING" && topup.ProviderPaymentUrl.Valid {
			return &model.VNPayPaymentURLResponse{PaymentURL: topup.ProviderPaymentUrl.String}, nil
		}
	}

    // 4. Sinh reference code an toàn
    randSuffix := uuid.New().String()[:8] // rút gọn cho gọn
    referenceCode := fmt.Sprintf("%s_%s_%s", time.Now().Format("060102"), order.ID.String(), randSuffix)

    // 5. Tạo topup order (status = pending)
    toArgs := database.CreateTopupOrderParams{
        UserID:        userInfo,
		OrderID:        utils.WrapUUID(order.ID),
        ReferenceCode: referenceCode,
        AmountBigint:  int64(order.TotalPrice), // đảm bảo order.TotalPrice là int64
        Currency:      "VND",
        Status:        "PENDING",
		Provider: 	"vnpay",
        IdempotencyKey: idempotencyKey,
    }
    topup, err = s.store.CreateTopupOrder(ctx, toArgs)
    if err != nil {
        // log chi tiết
        global.Logger.Sugar().Errorf("create topup fail: %v, user=%s order=%s", err, userInfo, order.ID)
        return nil, fmt.Errorf("failed to create topup order: %w", err)
    }

    // 6. Gọi service provider để tạo payment URL
    vnpaySvc := service.GetVNPayService()
    if vnpaySvc == nil {
        // mark topup failed
        _,err = s.store.UpdateTopupOrderStatus(ctx,database.UpdateTopupOrderStatusParams{
			 ID:     topup.ID,
			 Status: "failed",
		})
        return nil, errors.New("vnpay service not initialized")
    }

    paymentURL, err := vnpaySvc.GeneratePaymentURL(referenceCode, int(order.TotalPrice), ipAddr)
    if err != nil {
        // update topup status = failed để tránh dangling pending
        _,err = s.store.UpdateTopupOrderStatus(ctx, database.UpdateTopupOrderStatusParams{
			ID:     topup.ID,
			Status: "failed",
		})
        return nil, fmt.Errorf("failed to generate vnpay url: %w", err)
    }

    // 7. Update topup với paymentURL, expires, provider metadata
    _,err = s.store.UpdateTopupPaymentInfo(ctx, database.UpdateTopupPaymentInfoParams{
		ID:               topup.ID,
		ProviderPaymentUrl: sql.NullString{String: paymentURL, Valid: true},
		ExpiresAt:       sql.NullTime{Time: time.Now().Add(10 * time.Minute), Valid: true},
	})

    return &model.VNPayPaymentURLResponse{PaymentURL: paymentURL}, nil
}


func (s *paymentService) HandleVNPayCallback(ctx context.Context, params url.Values) error {
	vnpaySvc := service.GetVNPayService()
	if vnpaySvc == nil {
		return errors.New("vnpay service is not initialized")
	}
	receivedHash := params.Get("vnp_SecureHash")
	// 1. Xác thực chữ ký
	if !vnpaySvc.VerifySignature(params) {
		return errors.New("invalid vnpay signature")
	}

	// 2. Kiểm tra trạng thái giao dịch từ VNPay
	responseCode := params.Get("vnp_ResponseCode")
	if responseCode != "00" {
		return fmt.Errorf("vnpay transaction failed with response code: %s", responseCode)
	}

	// 3. Lấy thông tin TopupOrder từ DB
	refCode := params.Get("vnp_TxnRef")
	topupOrder, err := s.store.GetTopupOrderByReference(ctx, refCode)
	if err != nil {
		return fmt.Errorf("topup order with reference code %s not found: %w", refCode, err)
	}
	fmt.Println("status topup:", topupOrder.Status)
	// 4. Kiểm tra để tránh xử lý lại giao dịch đã hoàn thành
	if topupOrder.Status != "PENDING" {
		global.Logger.Sugar().Warnf("Received callback for an already processed topup order: %s, status: %s", topupOrder.ID, topupOrder.Status)
		return nil // Trả về nil để báo thành công cho VNPay, không xử lý lại
	}

	// 5. Kiểm tra số tiền
	vnpAmount, _ := strconv.ParseInt(params.Get("vnp_Amount"), 10, 64)
	if vnpAmount/100 != topupOrder.AmountBigint {
		return fmt.Errorf("amount mismatch: vnpay returned %d, expected %d", vnpAmount/100, topupOrder.AmountBigint)
	}

	// 6. Bắt đầu transaction để cập nhật CSDL
	err = s.store.ExecTx(ctx, func(q *database.Queries) error {
		// 6.1. Cập nhật trạng thái TopupOrder
		topupUpdated, err := q.UpdateTopupOrderStatus(ctx, database.UpdateTopupOrderStatusParams{
			ID:     topupOrder.ID,
			Status: "COMPLETED",
		})
		if err != nil {
			return fmt.Errorf("failed to update topup order status: %w", err)
		}

		// cập nhật order
		err = q.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
			ID:     topupUpdated.OrderID.UUID,
			Status: "paid",
		})
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		// 6.2. Ghi lại giao dịch thanh toán chi tiết
		payload, _ := json.Marshal(params)
		_, err = q.CreatePaymentTransaction(ctx, database.CreatePaymentTransactionParams{
			TopupOrderID:   utils.WrapUUID(topupOrder.ID),
			ProviderTxID:   params.Get("vnp_TransactionNo"),
			AmountBigint:   topupOrder.AmountBigint,
			Status:         "COMPLETED",
			Column6:  		payload,
			Signature:      sql.NullString{String:receivedHash, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create payment transaction: %w", err)
		}

		// 6.3. Lấy hoặc tạo ví cho người dùng
		wallet, err := q.GetWalletByUserID(ctx, topupOrder.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows)  {
				// Ví chưa tồn tại, tạo mới
				walletNew, err := q.CreateWallet(ctx, database.CreateWalletParams{
					UserID:        topupOrder.UserID,
					BalanceBigint: topupOrder.AmountBigint, // Khởi tạo ví với số dư 0
				})
				if err != nil {
					return fmt.Errorf("failed to create wallet for user %s: %w", topupOrder.UserID, err)
				}
				wallet = database.GetWalletByUserIDRow{
					ID: walletNew.ID,
					UserID: walletNew.UserID,
					BalanceBigint: walletNew.BalanceBigint,
				}
			} else {
				return fmt.Errorf("failed to get wallet for user %s: %w", topupOrder.UserID, err)
			}
		}

		// 6.4. Cộng tiền vào ví
		_, err = q.UpdateWalletBalance(ctx, database.UpdateWalletBalanceParams{
			ID:           wallet.ID,
			BalanceBigint: wallet.BalanceBigint + topupOrder.AmountBigint,
		})
		if err != nil {
			return fmt.Errorf("failed to credit wallet %s: %w", wallet.ID, err)
		}

		return nil
	})

	return err
}