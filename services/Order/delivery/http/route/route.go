package route

import (
	"github.com/gorilla/mux"
	"go-micro-services/services/Product/delivery/http"
)

type RootRoute struct {
	Router       *mux.Router
	ProductRoute *ProductRoute
}

func NewRootRoute(
	router *mux.Router,
	userRoute *ProductRoute,

) *RootRoute {
	rootRoute := &RootRoute{
		Router:       router,
		ProductRoute: userRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.ProductRoute.Register()
}

type ProductRoute struct {
	Router            *mux.Router
	ProductController *http.ProductController
}

func NewProductRoute(router *mux.Router, productController *http.ProductController) *ProductRoute {
	productRoute := &ProductRoute{
		Router:            router.PathPrefix("/products").Subrouter(),
		ProductController: productController,
	}
	return productRoute
}

func (productRoute *ProductRoute) Register() {
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.GetOneById).Methods("GET")
	productRoute.Router.HandleFunc("/update-stock/{id}", productRoute.ProductController.PatchOneById).Methods("PATCH")
}
