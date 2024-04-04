package container

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-micro-services/internal/config"
	httpdelivery "go-micro-services/internal/delivery/http"
	"go-micro-services/internal/delivery/http/route"
	"go-micro-services/internal/repository"
	"go-micro-services/internal/use_case"
)

type WebContainer struct {
	Env             *config.EnvConfig
	UserDatabase    *config.DatabaseConfig
	ProductDatabase *config.DatabaseConfig
	OrderDatabase   *config.DatabaseConfig
	Repository      *RepositoryContainer
	UseCase         *UseCaseContainer
	Controller      *ControllerContainer
	Route           *route.RootRoute
}

func NewWebContainer() *WebContainer {
	errEnvLoad := godotenv.Load()
	if errEnvLoad != nil {
		panic(fmt.Errorf("error loading .env file: %w", errEnvLoad))
	}

	envConfig := config.NewEnvConfig()
	userDBConfig := config.NewUserDBConfig(envConfig)
	productDBConfig := config.NewProductDBConfig(envConfig)
	orderDBConfig := config.NewOrderDBConfig(envConfig)

	userRepository := repository.NewUserRepository()
	productRepository := repository.NewProductRepository()
	orderRepository := repository.NewOrderRepository()
	repositoryContainer := NewRepositoryContainer(userRepository, productRepository, orderRepository)

	userUseCase := use_case.NewUserUseCase(userDBConfig, userRepository)
	productUseCase := use_case.NewProductUseCase(productDBConfig, productRepository)
	orderUseCase := use_case.NewOrderUseCase(orderDBConfig, orderRepository, envConfig)
	useCaseContainer := NewUseCaseContainer(userUseCase, productUseCase, orderUseCase)

	userController := httpdelivery.NewUserController(userUseCase)
	productController := httpdelivery.NewProductController(productUseCase)
	orderController := httpdelivery.NewOrderController(orderUseCase)
	controllerContainer := NewControllerContainer(userController, productController, orderController)

	router := mux.NewRouter()
	userRoute := route.NewUserRoute(router, userController)
	productRoute := route.NewProductRoute(router, productController)
	orderRoute := route.NewOrderRoute(router, orderController)

	rootRoute := route.NewRootRoute(
		router,
		userRoute,
		productRoute,
		orderRoute,
	)

	rootRoute.Register()

	webContainer := &WebContainer{
		Env:             envConfig,
		UserDatabase:    userDBConfig,
		ProductDatabase: productDBConfig,
		OrderDatabase:   orderDBConfig,
		Repository:      repositoryContainer,
		UseCase:         useCaseContainer,
		Controller:      controllerContainer,
		Route:           rootRoute,
	}

	return webContainer
}
