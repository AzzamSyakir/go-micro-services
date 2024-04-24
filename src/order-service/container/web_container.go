package container

import (
	"fmt"
	"go-micro-services/src/order-service/config"
	httpdelivery "go-micro-services/src/order-service/delivery/http"
	"go-micro-services/src/order-service/delivery/http/route"
	"go-micro-services/src/order-service/repository"
	"go-micro-services/src/order-service/use_case"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type WebContainer struct {
	Env           *config.EnvConfig
	OrderDatabase *config.DatabaseConfig
	Repository    *RepositoryContainer
	UseCase       *UseCaseContainer
	Controller    *ControllerContainer
	Route         *route.RootRoute
}

func NewWebContainer() *WebContainer {
	errEnvLoad := godotenv.Load()
	if errEnvLoad != nil {
		panic(fmt.Errorf("error loading .env file: %w", errEnvLoad))
	}

	envConfig := config.NewEnvConfig()
	orderDBConfig := config.NewDBConfig(envConfig)

	orderRepository := repository.NewOrderRepository()
	repositoryContainer := NewRepositoryContainer(orderRepository)

	orderUseCase := use_case.NewOrderUseCase(orderDBConfig, orderRepository, envConfig)

	useCaseContainer := NewUseCaseContainer(orderUseCase)

	orderController := httpdelivery.NewOrderController(orderUseCase)

	controllerContainer := NewControllerContainer(orderController)

	router := mux.NewRouter()
	orderRoute := route.NewOrderRoute(router, orderController)

	rootRoute := route.NewRootRoute(
		router,
		orderRoute,
	)

	rootRoute.Register()

	webContainer := &WebContainer{
		Env:           envConfig,
		OrderDatabase: orderDBConfig,
		Repository:    repositoryContainer,
		UseCase:       useCaseContainer,
		Controller:    controllerContainer,
		Route:         rootRoute,
	}

	return webContainer
}
