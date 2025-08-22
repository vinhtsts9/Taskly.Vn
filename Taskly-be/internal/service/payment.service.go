package service

import (
	"context"

	"net/url"

	model "Taskly.com/m/internal/models"
	"github.com/google/uuid"
)

type IPaymentService interface {
	// CreatePaymentIntent tạo bản ghi thanh toán và khởi tạo giao dịch với nhà cung cấp.
	CreatePaymentIntent(ctx context.Context, params model.CreatePaymentParams, userInfo uuid.UUID, ipAddr string,
    idempotencyKey string ) (*model.VNPayPaymentURLResponse, error)

	// HandleVNPayCallback xử lý IPN từ VNPay.
	HandleVNPayCallback(ctx context.Context, callbackData url.Values) error
}

var (
	localPaymentService IPaymentService
)

func GetPaymentService() IPaymentService {
	if localVNPayService == nil {
		panic("implement localVNPayService not found")
	}
	return localPaymentService
}

func InitPaymentService(s IPaymentService) {
	localPaymentService = s
}
