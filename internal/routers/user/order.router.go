package user

import (
	order_controller "Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (r *OrderRouter) InitOrderRouter(Router *gin.RouterGroup) {
	orderController := order_controller.NewOrderController()


	private := Router.Group("/orders")
	private.Use(middleware.AuthenMiddleware())
	{
		private.POST("/create", orderController.CreateOrder)
		private.GET("/list", orderController.ListOrdersByUser) // Route mới
		private.GET("/:id", orderController.GetOrderByID)
		private.PUT("/status", orderController.UpdateOrderStatus)
		private.PUT("/submit/:id", orderController.SubmitOrderDelivery)
		private.PUT("/complete/:id", orderController.AcceptOrderCompletion)
	}
}
