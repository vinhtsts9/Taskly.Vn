package initialize

import (
	"Taskly.com/m/global"
	"Taskly.com/m/internal/database"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.NewStore(global.PostgreSQL)
	service.InitUserService(impl.NewUserService(queries))
	service.InitGigService(impl.NewGigService(queries))
	service.InitOrderService(impl.NewOrderService(queries))
	service.InitChatService(impl.NewChatService(queries))
	service.InitDisputeService(impl.NewDisputeService(queries))
	service.InitVNPayService(impl.NewVNPayService(global.ENVSetting.Vnp_TmnCode, global.ENVSetting.Vnp_HashSecret, global.ENVSetting.Vnp_Url, global.ENVSetting.Vnp_UrlCallBack))
	service.InitRBACService(impl.NewRBACService(queries))
	service.InitAdminUserService(impl.NewAdminUserService(queries))
	service.InitPaymentService(impl.NewPaymentService(queries))
}
