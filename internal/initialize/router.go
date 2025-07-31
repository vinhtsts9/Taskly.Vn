package initialize

import (
	"Taskly.com/m/global"
	"Taskly.com/m/internal/middlewares"
	"Taskly.com/m/internal/routers"

	elasticsearch "Taskly.com/m/elasticSearch"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()

	}
	// middlewares
	// r.Use() logging
	// r.Use() cors
	r.Use(middlewares.CORSMiddleware())
	// r.Use() limiter global
	manageRouter := routers.RouterGroupApp.Manage
	userRouter := routers.RouterGroupApp.User

	MainGroup := r.Group("/v1/2024")
	{
		MainGroup.GET("/search", elasticsearch.SearchGigs)
		// MainGroup.GET("/checkStatus") //checking monitor

	}
	{
		userRouter.InitUserRouter(MainGroup)
		userRouter.InitOrderRouter(MainGroup)
		userRouter.InitGigRouter(MainGroup)
		userRouter.InitDisputeRouter(MainGroup)
		userRouter.InitChatRouter(MainGroup)
		userRouter.InitRbacRouter(MainGroup)
	}
	{
		manageRouter.InitAdminRouter(MainGroup)
		manageRouter.InitUserRouter(MainGroup)
	}
	return r
}
