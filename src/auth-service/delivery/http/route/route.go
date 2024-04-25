package route

import (
	"go-micro-services/src/auth-service/delivery/http"

	"github.com/gorilla/mux"
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
	UserController *http.AuthController
}

func NewUserRoute(router *mux.Router, userController *http.AuthController) *UserRoute {
	userRoute := &UserRoute{
		Router:         router.PathPrefix("/auth").Subrouter(),
		UserController: userController,
	}
	return userRoute
}

func (userRoute *UserRoute) Register() {
	userRoute.Router.HandleFunc("/login", userRoute.UserController.Login).Methods("POST")
}
