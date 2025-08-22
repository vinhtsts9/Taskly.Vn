package user

import (
	payment_controller "Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type PaymentRouter struct{}

func (r *PaymentRouter) InitPaymentRouter(Router *gin.RouterGroup) {
	paymentController := payment_controller.NewPaymentController()

	private := Router.Group("/payments")
	private.Use(middleware.AuthenMiddleware())
	{
		// Route để tạo yêu cầu thanh toán và lấy URL từ nhà cung cấp
		private.POST("/create-intent", paymentController.CreatePaymentIntent)
	}

	// Route công khai để VNPay gửi callback (IPN)
	public := Router.Group("/payments")
	{
		public.GET("/vnpay-callback", paymentController.HandleVNPayCallback)
	}
}