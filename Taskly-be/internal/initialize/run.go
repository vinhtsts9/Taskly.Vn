package initialize

import (
	"Taskly.com/m/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() *gin.Engine {
	LoadConfigRender()
	InitLogger()
	global.Logger.Info("Config ok", zap.String("ok", "success"))
	InitPostgreSQL()
	//InitCasbin()
	InitRedisFromEnvString()
	//InitKafka()
	NewCloudinary()
	//InitElasticSearch()

	InitServiceInterface()
	r := InitRouter()
	return r
}
