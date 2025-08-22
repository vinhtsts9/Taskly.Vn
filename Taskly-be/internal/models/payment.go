package model

import (
	"time"

	"github.com/google/uuid"
)

// ==== Core Payment Structs ====

// Payment represents a payment record in the database, linked to an order.
type Payment struct {
	ID                    uuid.UUID `json:"id"`
	OrderID               uuid.UUID `json:"order_id"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"` // e.g., "VND"
	Status                string    `json:"status"`   // e.g., "pending", "completed", "failed"
	PaymentMethod         string    `json:"payment_method"` // e.g., "zalopay", "vnpay", "credit_card"
	ProviderTransactionID *string   `json:"provider_transaction_id,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// PaymentTransaction represents an individual transaction attempt or status update.
type PaymentTransaction struct {
	ID                uuid.UUID `json:"id"`
	PaymentID         uuid.UUID `json:"payment_id"`
	Status            string    `json:"status"`
	ProviderResponse  *string   `json:"provider_response,omitempty"` // Raw response from payment provider
	CreatedAt         time.Time `json:"created_at"`
}

// CreatePaymentParams defines the input for creating a new payment intent.
type CreatePaymentParams struct {
	OrderID       uuid.UUID `json:"order_id" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required,oneof=zalopay vnpay"`
}

// ==== ZaloPay Integration Structs ====

// ZaloPayCreateOrderRequest is the structure for the request sent to ZaloPay to create a payment order.
type ZaloPayCreateOrderRequest struct {
	AppID       int64  `json:"app_id"`
	AppTransID  string `json:"app_trans_id"` // Format: yymmdd_orderid
	AppUser     string `json:"app_user"`
	AppTime     int64  `json:"app_time"` // Milliseconds
	Amount      int64  `json:"amount"`
	Item        string `json:"item"`       // JSON string of items
	Description string `json:"description"`
	EmbedData   string `json:"embed_data"` // JSON string for callback info
	BankCode    string `json:"bank_code"`  // "zalopayapp" or specific bank codes
	CallbackURL string `json:"callback_url"`
	Mac         string `json:"mac"`
}

// ZaloPayCreateOrderResponse is the structure of the response from ZaloPay after creating an order.
type ZaloPayCreateOrderResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	OrderURL      string `json:"order_url"`      // URL for user to proceed with payment
	ZpTransToken  string `json:"zp_trans_token"` // Transaction token
}

// ZaloPayCallbackRequest is the structure of the callback (Instant Payment Notification) from ZaloPay.
// ZaloPay sends this as a POST request with a JSON body.
type ZaloPayCallbackRequest struct {
	Data string `json:"data"` // A JSON string containing the transaction details
	Mac  string `json:"mac"`
}

// ZaloPayCallbackData is the structure of the JSON string inside the 'data' field of the callback.
type ZaloPayCallbackData struct {
	AppID          int64  `json:"app_id"`
	AppTransID     string `json:"app_trans_id"`
	AppTime        int64  `json:"app_time"`
	AppUser        string `json:"app_user"`
	Amount         int64  `json:"amount"`
	EmbedData      string `json:"embed_data"`
	Item           string `json:"item"`
	ZpTransID      int64  `json:"zp_trans_id"`
	ServerTime     int64  `json:"server_time"`
	Channel        int    `json:"channel"`
	MerchantUserID string `json:"merchant_user_id"`
	UserFeeAmount  int64  `json:"user_fee_amount"`
	DiscountAmount int64  `json:"discount_amount"`
}

// ==== VNPay Integration Structs ====

// VNPayPaymentURLResponse is the structure for the payment URL returned to the frontend.
type VNPayPaymentURLResponse struct {
	PaymentURL string `json:"payment_url"`
}

// VNPayCallbackParams represents the query parameters sent by VNPay in the return/IPN URL.
type VNPayCallbackParams struct {
	Amount            string `form:"vnp_Amount"`
	BankCode          string `form:"vnp_BankCode"`
	BankTranNo        string `form:"vnp_BankTranNo"`
	CardType          string `form:"vnp_CardType"`
	OrderInfo         string `form:"vnp_OrderInfo"`
	PayDate           string `form:"vnp_PayDate"`
	ResponseCode      string `form:"vnp_ResponseCode"`
	TmnCode           string `form:"vnp_TmnCode"`
	TransactionNo     string `form:"vnp_TransactionNo"`
	TransactionStatus string `form:"vnp_TransactionStatus"`
	TxnRef            string `form:"vnp_TxnRef"`
	SecureHash        string `form:"vnp_SecureHash"`
}
