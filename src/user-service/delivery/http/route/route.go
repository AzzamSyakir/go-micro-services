package route

import (
	"go-micro-services/src/user-service/delivery/http"

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
	userRoute.Router.HandleFunc("", userRoute.UserController.CreateUser).Methods("POST")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.GetOneById).Methods("GET")
	userRoute.Router.HandleFunc("/email/{email}", userRoute.UserController.GetOneByEmail).Methods("GET")
	userRoute.Router.HandleFunc("", userRoute.UserController.ListUser).Methods("GET")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.UpdateUser).Methods("PATCH")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.DeleteUser).Methods("DELETE")
}
