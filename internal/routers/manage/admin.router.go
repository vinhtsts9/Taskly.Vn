package manage

import "github.com/gin-gonic/gin"

type AdminRouter struct{}

func (ar *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	adminRouterPublic := Router.Group("/admin")
	{
		adminRouterPublic.POST("/login")
	}

}
