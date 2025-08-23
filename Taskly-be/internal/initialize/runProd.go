package initialize

import (
	"Taskly.com/m/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RunProd() *gin.Engine {
	LoadConfigProd()
	InitLogger()
	global.Logger.Info("Config ok", zap.String("ok", "success"))
	InitPostgreSQLProd()
	//InitCasbin()
	InitRedisProduction()
	//InitKafka()
	NewCloudinary()
	//InitElasticSearch()

	InitServiceInterface()
	r := InitRouter()
	return r
}
