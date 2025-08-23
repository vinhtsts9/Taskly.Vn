package initialize

import (
	"Taskly.com/m/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RunDev() *gin.Engine {
	LoadConfigDev()
	InitLogger()
	global.Logger.Info("Config ok", zap.String("ok", "success"))
	InitPostgreSQLDev()
	//InitCasbin()
	InitRedisDev()
	//InitKafka()
	NewCloudinary()
	//InitElasticSearch()

	InitServiceInterface()
	r := InitRouter()
	return r
}
