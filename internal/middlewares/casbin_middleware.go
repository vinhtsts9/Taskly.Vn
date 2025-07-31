package middlewares

import (
	"net/http"

	"Taskly.com/m/global"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils/auth"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo := auth.GetUserFromContext(c)

		userId := userInfo.ID
		sub := userId.String() // Chuyển UUID thành string cho Casbin
		obj := c.Request.URL.Path
		act := c.Request.Method

		global.Logger.Sugar().Infof("Checking access for user: %s, resource: %s, action: %s", sub, obj, act)

		ok, err := global.Casbin.Enforce(sub, obj, act)
		if err != nil {
			global.Logger.Sugar().Error("Error enforcing policy: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error 1"})
			c.Abort()
			return
		}
		policies, _ := global.Casbin.GetPolicy()
		global.Logger.Sugar().Info("Current policies: ", policies)

		if !ok {
			global.Logger.Sugar().Infof("Policy not found in cache, querying database for user: %s", sub)

			// Sửa: Truyền thẳng userId (uuid.UUID) vào hàm
			permissions, dbErr := service.GetRbacService().GetPermissionsByUserID(c, userId)
			if dbErr != nil || len(permissions) == 0 {
				global.Logger.Sugar().Warnf("No permissions found for user: %s", sub)
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
			global.Logger.Sugar().Info("Permissions found: ", permissions)
			// Thêm quyền mới vào Casbin enforcer
			for _, perm := range permissions {
				global.Casbin.AddPolicy(sub, perm.Resource, perm.Action)
			}

			// Ghi lại chính sách vào file policy.csv
			if err := global.Casbin.SavePolicy(); err != nil {
				global.Logger.Sugar().Errorf("Failed to save policy: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}

			// Kiểm tra lại quyền sau khi cập nhật chính sách
			ok, err = global.Casbin.Enforce(sub, obj, act)
			if err != nil || !ok {
				global.Logger.Sugar().Warnf("Access denied for user: %s", sub)
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
		}
		// Ghi log tất cả các chính sách hiện có

		c.Next()
	}
}
