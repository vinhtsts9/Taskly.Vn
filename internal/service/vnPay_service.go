package service

type IVNPayService interface {
	// Tạo URL thanh toán từ thông tin đơn hàng
	// GeneratePaymentURL(orderID string, amount int) (string, error)

	// // Xác minh chữ ký callback từ VNPay
	// VerifySignature(params url.Values) bool
}

var (
	localVNPayService IVNPayService
)

func GetVNPayService() IVNPayService {
	if localVNPayService == nil {
		panic("implement localVNPayService not found")
	}
	return localVNPayService
}

func InitVNPayService(s IVNPayService) {
	localVNPayService = s
}
