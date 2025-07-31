package user

import (
	"Taskly.com/m/internal/controller"
	"github.com/gin-gonic/gin"
)

type RbacRouter struct{}

func (r *RbacRouter) InitRbacRouter(rg *gin.RouterGroup) {
	controller := controller.NewRBACController() // inject service vào controller

	rbacRouter := rg.Group("/rbac")
	{
		rbacRouter.POST("/role", controller.CreateRole)
		rbacRouter.POST("/permission", controller.CreatePermission)
		rbacRouter.POST("/user-role", controller.AddRoleToUser)
		rbacRouter.POST("/role-permission", controller.AddPermissionToRole)
	}
}
