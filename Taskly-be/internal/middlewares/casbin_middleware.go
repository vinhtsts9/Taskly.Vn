package middlewares

import (
	"net/http"

	"Taskly.com/m/global"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy thông tin user từ context (đã được parse từ JWT)
		userInfo := auth.GetUserFromContext(c)
		if userInfo.ID == uuid.Nil { // Đảm bảo UserID có giá trị
			global.Logger.Sugar().Warn("User ID not found in context.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Ưu tiên dùng pattern route của Gin để khớp với permissions (có :param)
		resource := c.FullPath()
		if resource == "" {
			resource = c.Request.URL.Path
		}
		action := c.Request.Method

		// Gọi hàm CheckPermission tùy chỉnh của bạn
		hasPermission, err := service.GetRbacService().CheckPermission(c, userInfo.ID, resource, action)
		if err != nil {
			global.Logger.Sugar().Errorf("Lỗi khi kiểm tra quyền cho người dùng %s: %v", userInfo.ID.String(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if !hasPermission {
			global.Logger.Sugar().Warnf("Người dùng %s không có quyền truy cập %s %s", userInfo.ID.String(), action, resource)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		global.Logger.Sugar().Infof("Người dùng %s được cấp quyền truy cập %s %s", userInfo.ID.String(), action, resource)
		c.Next()
	}
}
