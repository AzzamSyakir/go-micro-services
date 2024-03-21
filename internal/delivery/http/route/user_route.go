package route

import (
	"github.com/gorilla/mux"
	"go-micro-services/internal/delivery/http"
)

type UserRoute struct {
	Router         *mux.Router
	UserController *http.UserController
}

func NewUserRoute(router *mux.Router, userController *http.UserController) *UserRoute {
	userRoute := &UserRoute{
		Router:         router.PathPrefix("/users").Subrouter(),
		UserController: userController,
	}
	return userRoute
}

func (userRoute *UserRoute) Register() {
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.GetById).Methods("GET")
}
