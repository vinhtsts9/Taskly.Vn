package user

type UserRouterGroup struct {
	UserRouter
	OrderRouter
	ChatRouter
	GigRouter
	DisputeRouter
	PaymentRouter
}
