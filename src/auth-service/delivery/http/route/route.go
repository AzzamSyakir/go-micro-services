package route

import (
	"go-micro-services/src/auth-service/delivery/http"
	"go-micro-services/src/auth-service/delivery/http/middleware"

	"github.com/gorilla/mux"
)

type AuthRoute struct {
	Middleware     *middleware.AuthMiddleware
	Router         *mux.Router
	AuthController *http.AuthController
}

func NewAuthRoute(router *mux.Router, AuthController *http.AuthController, middleware *middleware.AuthMiddleware) *AuthRoute {
	AuthRoute := &AuthRoute{
		Router:         router.PathPrefix("/auths").Subrouter(),
		AuthController: AuthController,
		Middleware:     middleware,
	}
	return AuthRoute
}

func (authRoute *AuthRoute) Register() {
	authRoute.Router.HandleFunc("/register", authRoute.AuthController.Register).Methods("POST")
	authRoute.Router.HandleFunc("/login", authRoute.AuthController.Login).Methods("POST")
	authRoute.Router.HandleFunc("/access-token", authRoute.AuthController.GetNewAccessToken).Methods("POST")
	authRoute.Router.HandleFunc("/logout", authRoute.AuthController.Logout).Methods("POST")
	protected := authRoute.Router.PathPrefix("").Subrouter()
	protected.Use(authRoute.Middleware.Middleware)
	protected.HandleFunc("/list-sessions", authRoute.AuthController.ListSession).Methods("GET")

}
