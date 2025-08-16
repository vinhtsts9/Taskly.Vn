package manage

type ManageRouterGroup struct {
	AdminRouter
}

func NewManageRouterGroup(adminRouter AdminRouter) ManageRouterGroup {
	return ManageRouterGroup{
		AdminRouter: adminRouter,
	}
}
