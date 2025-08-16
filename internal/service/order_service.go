package service

import (
	"context"

	model "Taskly.com/m/internal/models"
	"github.com/google/uuid"
)

type IOrderService interface {

	// 1. Tạo đơn hàng mới
	CreateOrder(ctx context.Context, input model.CreateOrderParams,buyerID uuid.UUID) (model.OrderResult, error)

	// 1b. Tạo đơn và sinh URL thanh toán VNPAY
	// CreateOrderAndGenerateVNPayURL(ctx context.Context, input model.CreateOrderParams, vnpayService IVNPayService) (string, error)

	// 1c. Xử lý callback từ VNPAY
	// HandleVNPayCallback(ctx context.Context, params url.Values, vnpayService IVNPayService) error

	// 2. Lấy danh sách đơn hàng của user
	ListOrdersByUser(ctx context.Context, userID uuid.UUID) ([]model.OrderResult, error)

	// 3. Lấy chi tiết đơn hàng
	GetOrderByID(ctx context.Context, id uuid.UUID) (model.OrderResult, error)

	// 4. Cập nhật trạng thái đơn hàng bất kỳ
	UpdateOrderStatus(ctx context.Context, input model.UpdateOrderStatusParams) error

	// 5. Freelancer giao hàng
	SubmitOrderDelivery(ctx context.Context, id uuid.UUID) error

	// 6. Buyer xác nhận hoàn thành
	AcceptOrderCompletion(ctx context.Context, id uuid.UUID) error
}

var (
	localOrderService IOrderService
)

func GetOrderService() IOrderService {
	if localOrderService == nil {
		panic("implement localOrderService not found")
	}
	return localOrderService
}

func InitOrderService(i IOrderService) {
	localOrderService = i
}
