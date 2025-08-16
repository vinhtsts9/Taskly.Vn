package manage

import (
	"Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type AdminRouter struct{}

func (ar *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	// Protected RBAC routes under /admin/rbac
	rbac := Router.Group("/admin/rbac")
	rbac.Use(middleware.AuthenMiddleware(), middleware.CasbinMiddleware())
	{
		ctrl := controller.NewRBACController()

		// Roles
		rbac.POST("/role", ctrl.CreateRole)
		// Optional list/delete có thể bổ sung sau

		// Permissions
		rbac.POST("/permission", ctrl.CreatePermission)

		// Role-Permissions mapping
		rbac.POST("/role-permission", ctrl.AddPermissionToRole)
		rbac.DELETE("/role-permission", ctrl.RemovePermissionFromRole)

		// User-Roles mapping
		rbac.POST("/user-role", ctrl.AddRoleToUser)
		rbac.DELETE("/user-role", ctrl.RemoveRoleFromUser)

		// Query helpers
		rbac.GET("/roles/:user_id/permissions", ctrl.GetPermissionsByUserID)
		rbac.GET("/roles/:user_id", ctrl.GetRolesByUserID)
		rbac.GET("/permissions/:role_id", ctrl.GetPermissionsByRoleID)
	}

	// Admin users listing
	admin := Router.Group("/admin")
	admin.Use(middleware.AuthenMiddleware(), middleware.CasbinMiddleware())
	{
		ctrl := controller.NewAdminUserController()
		admin.GET("/users", ctrl.ListUsers)
	}
}
