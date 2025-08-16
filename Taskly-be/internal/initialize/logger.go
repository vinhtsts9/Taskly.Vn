package initialize

import (
	"Taskly.com/m/global"
	"Taskly.com/m/package/logger"
)

func InitLogger() {
	lz := logger.NewLogger(global.Config.Logger)
	global.Logger = lz.GetZapLogger()
}
