package user

import (
	gig_controller "Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type GigRouter struct{}

func (r *GigRouter) InitGigRouter(Router *gin.RouterGroup) {
	gigController := gig_controller.NewGigController()

	public := Router.Group("/gigs")
	{
		public.GET("/:id", gigController.GetServiceByID)
		public.GET("/categories", gigController.GetCategories)
		public.GET("/search", gigController.SearchGigs)

	}

	private := Router.Group("/gigs")
	private.Use(middleware.AuthenMiddleware())
	{
		private.POST("", gigController.CreateService)
		private.PUT("", gigController.UpdateService)
		private.DELETE("/:id", gigController.DeleteService)

		// Upload gig media
		private.POST("upload-media", gigController.UploadGigMedia)
	}
}
