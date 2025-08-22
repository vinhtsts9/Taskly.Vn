package initialize

import (
	"Taskly.com/m/global"
	"Taskly.com/m/package/logger"
)

func InitLogger() {
	lz := logger.NewLogger(global.ENVSetting)
	global.Logger = lz.GetZapLogger()
}
