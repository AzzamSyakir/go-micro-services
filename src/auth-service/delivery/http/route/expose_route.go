package route

import (
	"go-micro-services/src/auth-service/delivery/http"
	"go-micro-services/src/auth-service/delivery/http/middleware"

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
	Middleware         *middleware.AuthMiddleware
	Router             *mux.Router
	CategoryController *http.ExposeController
}

func NewCategoryRoute(router *mux.Router, CategoryController *http.ExposeController, middleware *middleware.AuthMiddleware) *CategoryRoute {
	CategoryRoute := &CategoryRoute{
		Router:             router.PathPrefix("/categories").Subrouter(),
		CategoryController: CategoryController,
		Middleware:         middleware,
	}
	return CategoryRoute
}

func (CategoryRoute *CategoryRoute) Register() {
	CategoryRoute.Router.Use(CategoryRoute.Middleware.Middleware)
	CategoryRoute.Router.HandleFunc("", CategoryRoute.CategoryController.CreateCategory).Methods("POST")
	CategoryRoute.Router.HandleFunc("", CategoryRoute.CategoryController.ListCategories).Methods("GET")
	CategoryRoute.Router.HandleFunc("/{id}", CategoryRoute.CategoryController.DetailCategory).Methods("GET")
	CategoryRoute.Router.HandleFunc("/{id}", CategoryRoute.CategoryController.DeleteCategory).Methods("DELETE")
	CategoryRoute.Router.HandleFunc("/{id}", CategoryRoute.CategoryController.UpdateCategory).Methods("PATCH")
}

// order route

type OrderRoute struct {
	Middleware      *middleware.AuthMiddleware
	Router          *mux.Router
	OrderController *http.ExposeController
}

func NewOrderRoute(router *mux.Router, orderController *http.ExposeController, middleware *middleware.AuthMiddleware) *OrderRoute {
	orderRoute := &OrderRoute{
		Router:          router.PathPrefix("/orders").Subrouter(),
		OrderController: orderController,
		Middleware:      middleware,
	}
	return orderRoute
}
func (orderRoute *OrderRoute) Register() {
	orderRoute.Router.Use(orderRoute.Middleware.Middleware)
	orderRoute.Router.HandleFunc("", orderRoute.OrderController.Orders).Methods("POST")
}

// product route

type ProductRoute struct {
	Middleware        *middleware.AuthMiddleware
	Router            *mux.Router
	ProductController *http.ExposeController
}

func NewProductRoute(router *mux.Router, productController *http.ExposeController, middleware *middleware.AuthMiddleware) *ProductRoute {
	productRoute := &ProductRoute{
		Router:            router.PathPrefix("/products").Subrouter(),
		ProductController: productController,
		Middleware:        middleware,
	}
	return productRoute
}

func (productRoute *ProductRoute) Register() {
	productRoute.Router.Use(productRoute.Middleware.Middleware)
	productRoute.Router.HandleFunc("", productRoute.ProductController.CreateProduct).Methods("POST")
	productRoute.Router.HandleFunc("", productRoute.ProductController.ListProducts).Methods("GET")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.DetailProduct).Methods("GET")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.DeleteProduct).Methods("DELETE")
	productRoute.Router.HandleFunc("/{id}", productRoute.ProductController.UpdateProduct).Methods("PATCH")
}

// user route

type UserRoute struct {
	Middleware     *middleware.AuthMiddleware
	Router         *mux.Router
	UserController *http.ExposeController
}

func NewUserRoute(router *mux.Router, userController *http.ExposeController, middleware *middleware.AuthMiddleware) *UserRoute {
	userRoute := &UserRoute{
		Router:         router.PathPrefix("/users").Subrouter(),
		UserController: userController,
		Middleware:     middleware,
	}
	return userRoute
}

func (userRoute *UserRoute) Register() {
	userRoute.Router.Use(userRoute.Middleware.Middleware)
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.DetailUser).Methods("GET")
	userRoute.Router.HandleFunc("/email/{email}", userRoute.UserController.GetUserByEmail).Methods("GET")
	userRoute.Router.HandleFunc("", userRoute.UserController.FetchUser).Methods("GET")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.UpdateUser).Methods("PATCH")
	userRoute.Router.HandleFunc("/{id}", userRoute.UserController.DeleteUser).Methods("DELETE")
}
