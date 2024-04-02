package route

import (
	"github.com/gorilla/mux"
	"go-micro-services/internal/delivery/http"
)

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
	productRoute.Router.HandleFunc("", productRoute.OrderController.Orders).Methods("POST")
}
