package user

import (
	"Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userController := controller.NewUserController()

	userRouterPublic := Router.Group("/users")
	{
		userRouterPublic.POST("/register", userController.Register)
		userRouterPublic.POST("/verify-otp", userController.VerifyOTP)
		userRouterPublic.POST("/update-password-register", userController.UpdatePasswordRegister)
		userRouterPublic.POST("/login", userController.Login)
		userRouterPublic.POST("/refresh-token", userController.RefreshToken) // Route mới
		userRouterPublic.GET("/:id", userController.GetUserByID)
	}

	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(middleware.AuthenMiddleware())
	{
		userRouterPrivate.GET("/me", userController.Me) // Route mới để kiểm tra phiên
		userRouterPrivate.POST("/logout", userController.Logout)
		userRouterPrivate.DELETE("/", userController.DeleteUser)
	}
}
