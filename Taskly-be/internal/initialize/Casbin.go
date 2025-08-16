package initialize

import (
	"Taskly.com/m/global"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func InitCasbin() {
	modelPath := "configs/model.conf"
	policyPath := "configs/policy.csv"

	adapter := fileadapter.NewAdapter(policyPath)

	e, err := casbin.NewEnforcer(modelPath, adapter)

	if err != nil {
		global.Logger.Sugar().Error("Failed to create enforcer %v:", err)
	}
	global.Casbin = e
	if global.Casbin != nil {
		global.Logger.Sugar().Info("Casbin enforcer initialization success")
	}
}
