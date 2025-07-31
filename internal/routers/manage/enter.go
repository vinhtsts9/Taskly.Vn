package manage

type ManageRouterGroup struct {
	UserRouter
	AdminRouter
}

func NewManageRouterGroup(userRouter UserRouter, adminRouter AdminRouter) ManageRouterGroup {
	return ManageRouterGroup{
		UserRouter:  userRouter,
		AdminRouter: adminRouter,
	}
}
