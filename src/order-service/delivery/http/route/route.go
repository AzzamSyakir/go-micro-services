package route

import (
	"go-micro-services/src/order-service/delivery/http"

	"github.com/gorilla/mux"
)

type RootRoute struct {
	Router     *mux.Router
	OrderRoute *OrderRoute
}

func NewRootRoute(
	router *mux.Router,
	userRoute *OrderRoute,

) *RootRoute {
	rootRoute := &RootRoute{
		Router:     router,
		OrderRoute: userRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.OrderRoute.Register()
}

type OrderRoute struct {
	Router          *mux.Router
	OrderController *http.OrderController
}

func NewOrderRoute(router *mux.Router, orderController *http.OrderController) *OrderRoute {
	orderRoute := &OrderRoute{
		Router:          router.PathPrefix("/orders").Subrouter(),
		OrderController: orderController,
	}
	return orderRoute
}

func (productRoute *OrderRoute) Register() {
	productRoute.Router.HandleFunc("/{id}", productRoute.OrderController.Orders).Methods("POST")
}
