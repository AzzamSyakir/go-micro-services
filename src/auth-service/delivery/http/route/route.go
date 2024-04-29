package route

import (
	"go-micro-services/src/auth-service/delivery/http"

	"github.com/gorilla/mux"
)

type RootRoute struct {
	Router    *mux.Router
	AuthRoute *AuthRoute
}

func NewRootRoute(
	router *mux.Router,
	authRoute *AuthRoute,

) *RootRoute {
	rootRoute := &RootRoute{
		Router:    router,
		AuthRoute: authRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.AuthRoute.Register()
}

type AuthRoute struct {
	Router         *mux.Router
	AuthController *http.AuthController
}

func NewAuthRoute(router *mux.Router, AuthController *http.AuthController) *AuthRoute {
	AuthRoute := &AuthRoute{
		Router:         router.PathPrefix("/auths").Subrouter(),
		AuthController: AuthController,
	}
	return AuthRoute
}

func (authRoute *AuthRoute) Register() {
	authRoute.Router.HandleFunc("/register", authRoute.AuthController.Register).Methods("POST")

	authRoute.Router.HandleFunc("/login", authRoute.AuthController.Login).Methods("POST")
	authRoute.Router.HandleFunc("/access-token", authRoute.AuthController.Login).Methods("POST")
	authRoute.Router.HandleFunc("/logout", authRoute.AuthController.Login).Methods("POST")
}
