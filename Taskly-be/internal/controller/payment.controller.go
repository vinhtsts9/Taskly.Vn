package controller

import (
	"net/http"
	"net/url"

	"Taskly.com/m/global"

	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils/auth"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	svc service.IPaymentService
}

func NewPaymentController() *PaymentController {
	return &PaymentController{
		svc: service.GetPaymentService(), // Giả sử có một hàm singleton
	}
}

// CreatePaymentIntent xử lý yêu cầu tạo thanh toán mới.
// Nó trả về URL thanh toán từ nhà cung cấp (ví dụ: ZaloPay).
func (ctl *PaymentController) CreatePaymentIntent(c *gin.Context) {
	var params model.CreatePaymentParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	idempotencyKey := c.GetHeader("Idempotency-Key")


	userInfo := auth.GetUserFromContext(c)
	if userInfo == nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    return
}

	paymentResponse, err := ctl.svc.CreatePaymentIntent(c, params, userInfo.ID, c.ClientIP(), idempotencyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	global.Logger.Sugar().Infof("Payment intent created: %+v", paymentResponse)
	c.JSON(http.StatusOK, gin.H{"payment_url": paymentResponse})
}

// HandleVNPayCallback xử lý thông báo thanh toán từ VNPay.
func (ctl *PaymentController) HandleVNPayCallback(c *gin.Context) {
       err := ctl.svc.HandleVNPayCallback(c.Request.Context(), c.Request.URL.Query())
       if err != nil {
	       global.Logger.Sugar().Errorf("Error handling VNPay callback: %v", err)
	       // Có thể redirect về FE với trạng thái thất bại
	       feURL := "https://taskly-vn.vercel.app/payment-result?success=false&message=" + url.QueryEscape(err.Error())
	       c.Redirect(http.StatusFound, feURL)
	       return
       }
       // Nếu thành công, redirect về FE với trạng thái thành công
       feURL := "https://taskly-vn.vercel.app/payment-result?success=true"
       c.Redirect(http.StatusFound, feURL)
}