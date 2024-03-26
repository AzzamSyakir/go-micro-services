package route

import (
	"github.com/gorilla/mux"
	"go-micro-services/internal/delivery/http"
)

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
