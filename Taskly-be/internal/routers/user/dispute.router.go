package user

import (
	dispute_controller "Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type DisputeRouter struct{}

func (r *DisputeRouter) InitDisputeRouter(Router *gin.RouterGroup) {
	disputeController := dispute_controller.NewDisputeController()

	private := Router.Group("/disputes")
	private.Use(middleware.AuthenMiddleware())
	{
		private.POST("/", disputeController.CreateDispute)
		private.GET("/", disputeController.ListDisputes)
		private.PUT("/status", disputeController.UpdateDisputeStatus)
	}
}
