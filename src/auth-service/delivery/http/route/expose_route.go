package route

import (
	"go-micro-services/src/auth-service/delivery/http"

	"github.com/gorilla/mux"
)

type ExposeRoute struct {
	Router        *mux.Router
	UserRoute     *UserRoute
	ProductRoute  *ProductRoute
	CategoryRoute *CategoryRoute
	OrderRoute    *OrderRoute
}

func NewExposeRoute(
	router *mux.Router,
	userRoute *UserRoute,
	productRoute *ProductRoute,
	categoryRoute *CategoryRoute,
	orderRoute *OrderRoute,

) *ExposeRoute {
	rootRoute := &ExposeRoute{
		Router:        router,
		UserRoute:     userRoute,
		ProductRoute:  productRoute,
		CategoryRoute: categoryRoute,
		OrderRoute:    orderRoute,
	}
	return rootRoute
}

func (exposeRoute *ExposeRoute) Register() {
	exposeRoute.UserRoute.Register()
	exposeRoute.ProductRoute.Register()
	exposeRoute.CategoryRoute.Register()
	exposeRoute.OrderRoute.Register()
}

type CategoryRoute struct {
	CategoryRouteRouter *mux.Router
	CategoryController  *http.ExposeController
}

func NewCategoryRoute(router *mux.Router, CategoryController *http.ExposeController) *CategoryRoute {
	CategoryRoute := &CategoryRoute{
		CategoryRouteRouter: router.PathPrefix("/categories").Subrouter(),
		CategoryController:  CategoryController,
	}
	return CategoryRoute
}

func (CategoryRoute *CategoryRoute) Register() {
	CategoryRoute.CategoryRouteRouter.HandleFunc("", CategoryRoute.CategoryController.CreateCategory).Methods("POST")
	CategoryRoute.CategoryRouteRouter.HandleFunc("", CategoryRoute.CategoryController.ListCategories).Methods("GET")
	CategoryRoute.CategoryRouteRouter.HandleFunc("/{id}", CategoryRoute.CategoryController.ListCategories).Methods("GET")
	CategoryRoute.CategoryRouteRouter.HandleFunc("/{id}", CategoryRoute.CategoryController.DeleteCategory).Methods("DELETE")
	CategoryRoute.CategoryRouteRouter.HandleFunc("/{id}", CategoryRoute.CategoryController.UpdateCategory).Methods("PATCH")
}

// order route

type OrderRoute struct {
	Router          *mux.Router
	OrderController *http.ExposeController
}

func NewOrderRoute(router *mux.Router, orderController *http.ExposeController) *OrderRoute {
	orderRoute := &OrderRoute{
		Router:          router.PathPrefix("/orders").Subrouter(),
		OrderController: orderController,
	}
	return orderRoute
}
func (productRoute *OrderRoute) Register() {
	productRoute.Router.HandleFunc("/{id}", productRoute.OrderController.Orders).Methods("POST")
}

// product route

type ProductRoute struct {
	Router            *mux.Router
	ProductController *http.ExposeController
}

func NewProductRoute(router *mux.Router, productController *http.ExposeController) *ProductRoute {
	productRoute := &ProductRoute{
		Router:            router.PathPrefix("/products").Subrouter(),
		ProductController: productController,
	}
	return productRoute
}

func (productRoute *ProductRoute) Register() {
	productRoute.Router.HandleFunc("", productRoute.ProductController.CreateProduct).Methods("POST")
	productRoute.Router.HandleFunc("", productRoute.ProductController.ListProducts).Methods("GET")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.DetailProduct).Methods("GET")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.DeleteProduct).Methods("DELETE")
	productRoute.Router.HandleFunc("/update-stock/{id}", productRoute.ProductController.UpdateStock).Methods("PATCH")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.UpdateProduct).Methods("PATCH")
}

// user route

type UserRoute struct {
	Router         *mux.Router
	UserController *http.ExposeController
}

func NewUserRoute(router *mux.Router, userController *http.ExposeController) *UserRoute {
	userRoute := &UserRoute{
		Router:         router.PathPrefix("/users").Subrouter(),
		UserController: userController,
	}
	return userRoute
}

func (userRoute *UserRoute) Register() {
	userRoute.Router.HandleFunc("", userRoute.UserController.CreateUser).Methods("POST")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.DetailUser).Methods("GET")
	userRoute.Router.HandleFunc("/email/{email}", userRoute.UserController.GetUserByEmail).Methods("GET")
	userRoute.Router.HandleFunc("", userRoute.UserController.FetchUser).Methods("GET")
	userRoute.Router.HandleFunc("/update-balance/{id}", userRoute.UserController.UpdateBalance).Methods("PATCH")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.UpdateUser).Methods("PATCH")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.DeleteUser).Methods("DELETE")
}
