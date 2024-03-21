package route

import (
	"github.com/gorilla/mux"
)

type RootRoute struct {
	Router       *mux.Router
	UserRoute    *UserRoute
	ProductRoute *ProductRoute
}

func NewRootRoute(
	router *mux.Router,
	userRoute *UserRoute,
	ProductRoute *ProductRoute,
) *RootRoute {
	rootRoute := &RootRoute{
		Router:       router,
		UserRoute:    userRoute,
		ProductRoute: ProductRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.UserRoute.Register()
	rootRoute.ProductRoute.Register()
}
