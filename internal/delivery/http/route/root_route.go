package route

import (
	"github.com/gorilla/mux"
)

type RootRoute struct {
	Router       *mux.Router
	UserRoute    *UserRoute
	ProductRoute *ProductRoute
	OrderRoute   *OrderRoute
}

func NewRootRoute(
	router *mux.Router,
	userRoute *UserRoute,
	productRoute *ProductRoute,
	orderRoute *OrderRoute,
) *RootRoute {
	rootRoute := &RootRoute{
		Router:       router,
		UserRoute:    userRoute,
		ProductRoute: productRoute,
		OrderRoute:   orderRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.UserRoute.Register()
	rootRoute.ProductRoute.Register()
	rootRoute.OrderRoute.Register()
}
