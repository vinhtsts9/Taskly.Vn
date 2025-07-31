package user

import (
	order_controller "Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (r *OrderRouter) InitOrderRouter(Router *gin.RouterGroup) {
	orderController := order_controller.NewOrderController()

	public := Router.Group("/orders")
	{
		public.POST("/", orderController.CreateOrder)
		public.POST("/vnpay", orderController.CreateOrderAndGetVNPayURL)
		public.GET("/vnpay/callback", orderController.HandleVNPayCallback)
		// Đã chuyển route lấy danh sách đơn hàng sang private group

	}

	private := Router.Group("/orders")
	private.Use(middleware.AuthenMiddleware())
	{
		private.GET("/", orderController.ListOrdersByUser) // Route mới
		private.GET("/:id", orderController.GetOrderByID)
		private.PUT("/status", orderController.UpdateOrderStatus)
		private.PUT("/submit/:id", orderController.SubmitOrderDelivery)
		private.PUT("/complete/:id", orderController.AcceptOrderCompletion)
	}
}
