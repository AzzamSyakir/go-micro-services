package route

import (
	"github.com/gorilla/mux"
	"go-micro-services/internal/delivery/http"
)

type ProductRoute struct {
	Router         *mux.Router
	UserController *http.UserController
}

func NewProductRoute(router *mux.Router, userController *http.UserController) *UserRoute {
	userRoute := &UserRoute{
		Router:         router.PathPrefix("/users").Subrouter(),
		UserController: userController,
	}
	return userRoute
}

func (userRoute *ProductRoute) Register() {
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.GetById).Methods("GET")
}
