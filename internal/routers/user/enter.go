package user

type UserRouterGroup struct {
	UserRouter
	OrderRouter
	RbacRouter
	ChatRouter
	GigRouter
	DisputeRouter
}
