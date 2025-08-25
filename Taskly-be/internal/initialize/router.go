package initialize

import (
	"Taskly.com/m/global"
	"Taskly.com/m/internal/routers"

	elasticsearch "Taskly.com/m/elasticSearch"
	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000",global.ENVSetting.Fe_api},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Idempotency-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))
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
		userRouter.InitPaymentRouter(MainGroup)
	}
	{
		manageRouter.InitAdminRouter(MainGroup)
	}
	return r
}
