package initialize

import (
	"Taskly.com/m/global"
	"Taskly.com/m/package/logger"
	"Taskly.com/m/package/setting"
)

func InitLogger(env *setting.ENV) {
	lz := logger.NewLogger(global.ENVSetting)
	global.Logger = lz.GetZapLogger()
}
