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
	AuthRoute *AuthRoute,

) *RootRoute {
	rootRoute := &RootRoute{
		Router:    router,
		AuthRoute: AuthRoute,
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
		Router:         router.PathPrefix("/auth").Subrouter(),
		AuthController: AuthController,
	}
	return AuthRoute
}

func (AuthRoute *AuthRoute) Register() {
	AuthRoute.Router.HandleFunc("/login", AuthRoute.AuthController.Login).Methods("POST")
}
