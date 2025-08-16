package routers

import (
	"Taskly.com/m/internal/routers/manage"
	"Taskly.com/m/internal/routers/user"
)

type RouterGroup struct {
	User   user.UserRouterGroup
	Manage manage.ManageRouterGroup
}

var RouterGroupApp = new(RouterGroup)
