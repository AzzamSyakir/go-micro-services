package route

import (
	"github.com/gorilla/mux"
	"go-micro-services/services/User/delivery/http"
)

type RootRoute struct {
	Router    *mux.Router
	UserRoute *UserRoute
}

func NewRootRoute(
	router *mux.Router,
	userRoute *UserRoute,

) *RootRoute {
	rootRoute := &RootRoute{
		Router:    router,
		UserRoute: userRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.UserRoute.Register()
}

type UserRoute struct {
	Router         *mux.Router
	UserController *http.UserController
}

func NewUserRoute(router *mux.Router, userController *http.UserController) *UserRoute {
	userRoute := &UserRoute{
		Router:         router.PathPrefix("/User").Subrouter(),
		UserController: userController,
	}
	return userRoute
}

func (userRoute *UserRoute) Register() {
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.GetOneById).Methods("GET")
	userRoute.Router.HandleFunc("/update-balance/{id}", userRoute.UserController.PatchOneById).Methods("PATCH")
}