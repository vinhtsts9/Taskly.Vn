package impl

import (
	"context"
	"encoding/json"

	"Taskly.com/m/global"
	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/package/utils"
	"Taskly.com/m/package/utils/mapper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
func (s *sOrderService) CreateOrder(ctx context.Context, input model.CreateOrderParams, buyerID uuid.UUID) (model.OrderResult, error) {
	var createdOrder database.Order

	// Bước 1: Lấy dữ liệu gốc từ Database bằng gig_id do client cung cấp.
	gigAndPackages, err := s.store.GetGigAndPackagesForOrder(ctx, input.GigID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.OrderResult{}, errors.New("gig not found")
		}
		global.Logger.Error("Failed to get gig with packages for order", zap.Error(err), zap.String("gigId", input.GigID.String()))
		return model.OrderResult{}, errors.Wrap(err, "failed to get gig details")
	}

	// Bước 2: Xác thực người bán. So sánh seller_id từ client với user_id (người bán thật) từ DB.
	if gigAndPackages.UserID != input.SellerID {
		return model.OrderResult{}, errors.New("seller ID mismatch")
	}

	// Giải nén thông tin các gói dịch vụ từ JSON
	var packages []struct {
		Tier  string  `json:"tier"`
		Price float64 `json:"price"`
	}
	if gigAndPackages.GigPackages != nil {
		if err := json.Unmarshal(gigAndPackages.GigPackages, &packages); err != nil {
			global.Logger.Error("Failed to unmarshal gig packages", zap.Error(err), zap.String("gigId", input.GigID.String()))
			return model.OrderResult{}, errors.Wrap(err, "failed to process gig packages")
		}
	}

	// Bước 3: Xác thực giá. Tìm giá thật trong dữ liệu từ DB.
	var serverPrice float64
	var foundPackage bool
	for _, pkg := range packages {
		if pkg.Tier == input.PackageTier {
			serverPrice = pkg.Price
			foundPackage = true
			break
		}
	}
	if !foundPackage {
		return model.OrderResult{}, errors.New("package tier not found for this gig")
	}

	err = s.store.ExecTx(ctx, func(q *database.Queries) error {
		var err error
		// Tạo Order với dữ liệu đã được xác thực
		createdOrder, err = q.CreateOrder(ctx, database.CreateOrderParams{
			GigID:        input.GigID,
			BuyerID:      buyerID, // Sử dụng ID của người dùng đã xác thực
			SellerID:     gigAndPackages.UserID, // Sử dụng ID người bán thật từ DB
			TotalPrice:   serverPrice, // Sử dụng giá từ server
			DeliveryDate: utils.ToNullTime(input.DeliveryDate),
		})
		if err != nil {
			return err
		}

		// 2. Lặp qua các câu trả lời và tạo chúng
		for _, answer := range input.Answers {
			_, err = q.CreateAnswer(ctx, database.CreateAnswerParams{
				GigID:      input.GigID,
				UserID:     buyerID, // Gắn câu trả lời với người mua
				QuestionID: answer.QuestionID,
				Answer:     answer.Answer,
			})
			if err != nil {
				return err // Rollback transaction
			}
		}
		return nil
	})

	if err != nil {
		return model.OrderResult{}, err
	}

	return mapper.ConvertDBOrderToModel(createdOrder), nil
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
