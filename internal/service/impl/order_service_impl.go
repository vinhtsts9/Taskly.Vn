package impl

import (
	"context"
	"fmt"
	"net/url"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils"
	"Taskly.com/m/package/utils/mapper"

	"github.com/google/uuid"
)

type sOrderService struct {
	store database.Store
}

func NewOrderService(store database.Store) *sOrderService {
	return &sOrderService{
		store: store,
	}
}

// 1. Tạo đơn hàng mới
func (s *sOrderService) CreateOrder(ctx context.Context, input model.CreateOrderParams) (model.OrderResult, error) {
	order, err := s.store.CreateOrder(ctx, database.CreateOrderParams{
		GigID:        input.GigID,
		BuyerID:      input.BuyerID,
		SellerID:     input.SellerID,
		TotalPrice:   input.TotalPrice,
		DeliveryDate: utils.ToNullTime(input.DeliveryDate),
	})
	if err != nil {
		return model.OrderResult{}, err
	}

	return mapper.ConvertDBOrderToModel(order), nil
}

func (s *sOrderService) CreateOrderAndGenerateVNPayURL(
	ctx context.Context,
	input model.CreateOrderParams,
	vnpayService service.IVNPayService,
) (string, error) {
	order, err := s.CreateOrder(ctx, input)
	if err != nil {
		return "", err
	}

	paymentURL, err := vnpayService.GeneratePaymentURL(order.ID.String(), int(order.TotalPrice*100))
	if err != nil {
		return "", err
	}

	return paymentURL, nil
}

func (s *sOrderService) HandleVNPayCallback(
	ctx context.Context,
	params url.Values,
	vnpayService service.IVNPayService,
) error {
	// 1. Verify chữ ký
	if !vnpayService.VerifySignature(params) {
		return fmt.Errorf("invalid signature")
	}

	// 2. Kiểm tra mã phản hồi
	if params.Get("vnp_ResponseCode") != "00" {
		return fmt.Errorf("payment failed")
	}

	// 3. Lấy orderID từ params và cập nhật trạng thái
	orderIDStr := params.Get("vnp_TxnRef")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return err
	}

	return s.UpdateOrderStatus(ctx, model.UpdateOrderStatusParams{
		ID:     orderID,
		Status: "paid",
	})
}

// 2. Lấy danh sách đơn hàng của user
func (s *sOrderService) ListOrdersByUser(ctx context.Context, userID uuid.UUID) ([]model.OrderResult, error) {
	orders, err := s.store.ListOrdersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapper.ConvertDBOrderListToModelList(orders), nil
}

// 3. Lấy chi tiết đơn hàng
func (s *sOrderService) GetOrderByID(ctx context.Context, id uuid.UUID) (model.OrderResult, error) {
	order, err := s.store.GetOrderByID(ctx, id)
	if err != nil {
		return model.OrderResult{}, err
	}

	return mapper.ConvertDBOrderToModel(order), nil
}

// 4. Cập nhật trạng thái bất kỳ
func (s *sOrderService) UpdateOrderStatus(ctx context.Context, input model.UpdateOrderStatusParams) error {
	return s.store.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     input.ID,
		Status: input.Status,
	})
}

// 5. Freelancer giao hàng
func (s *sOrderService) SubmitOrderDelivery(ctx context.Context, id uuid.UUID) error {
	return s.store.SubmitOrderDelivery(ctx, id)
}

// 6. Buyer xác nhận hoàn thành
func (s *sOrderService) AcceptOrderCompletion(ctx context.Context, id uuid.UUID) error {
	return s.store.AcceptOrderCompletion(ctx, id)
}
