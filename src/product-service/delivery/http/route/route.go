package route

import (
	"go-micro-services/src/product-service/delivery/http"

	"github.com/gorilla/mux"
)

type RootRoute struct {
	Router        *mux.Router
	ProductRoute  *ProductRoute
	CategoryRoute *CategoryRoute
}

func NewRootRoute(
	router *mux.Router,
	userRoute *ProductRoute,
	categoryRoute *CategoryRoute,

) *RootRoute {
	rootRoute := &RootRoute{
		Router:        router,
		ProductRoute:  userRoute,
		CategoryRoute: categoryRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.ProductRoute.Register()
	rootRoute.CategoryRoute.Register()
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
	productRoute.Router.HandleFunc("", productRoute.ProductController.CreateProduct).Methods("POST")
	productRoute.Router.HandleFunc("", productRoute.ProductController.ListProduct).Methods("GET")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.GetProduct).Methods("GET")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.DeleteProduct).Methods("DELETE")
	productRoute.Router.HandleFunc("/update-stock/{id}", productRoute.ProductController.UpdateStock).Methods("PATCH")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.UpdateProduct).Methods("PATCH")
}

type CategoryRoute struct {
	CategoryRouteRouter *mux.Router
	CategoryController  *http.CategoryController
}

func NewCategoryRoute(router *mux.Router, CategoryController *http.CategoryController) *CategoryRoute {
	CategoryRoute := &CategoryRoute{
		CategoryRouteRouter: router.PathPrefix("/categories").Subrouter(),
		CategoryController:  CategoryController,
	}
	return CategoryRoute
}

func (CategoryRoute *CategoryRoute) Register() {
	CategoryRoute.CategoryRouteRouter.HandleFunc("", CategoryRoute.CategoryController.CreateCategory).Methods("POST")
	CategoryRoute.CategoryRouteRouter.HandleFunc("", CategoryRoute.CategoryController.ListCategories).Methods("GET")
	CategoryRoute.CategoryRouteRouter.HandleFunc("/{id}", CategoryRoute.CategoryController.GetCategory).Methods("GET")
	CategoryRoute.CategoryRouteRouter.HandleFunc("/{id}", CategoryRoute.CategoryController.DeleteCategory).Methods("DELETE")
	CategoryRoute.CategoryRouteRouter.HandleFunc("/{id}", CategoryRoute.CategoryController.UpdateCategory).Methods("PATCH")
}
