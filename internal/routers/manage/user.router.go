package manage

import "github.com/gin-gonic/gin"

type UserRouter struct{}

func (ur *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouterPrivate := Router.Group("/admin/user")
	{
		userRouterPrivate.POST("/active_user")
	}
}
