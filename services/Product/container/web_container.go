package container

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-micro-services/services/Product/config"
	httpdelivery "go-micro-services/services/Product/delivery/http"
	"go-micro-services/services/Product/delivery/http/route"
	"go-micro-services/services/Product/repository"
	"go-micro-services/services/Product/use_case"
)

type WebContainer struct {
	Env             *config.EnvConfig
	ProductDatabase *config.DatabaseConfig
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
	productDBConfig := config.NewDBConfig(envConfig)

	productRepository := repository.NewProductRepository()
	repositoryContainer := NewRepositoryContainer(productRepository)

	productUseCase := use_case.NewProductUseCase(productDBConfig, productRepository)

	useCaseContainer := NewUseCaseContainer(productUseCase)

	productController := httpdelivery.NewProductController(productUseCase)

	controllerContainer := NewControllerContainer(productController)

	router := mux.NewRouter()
	productRoute := route.NewProductRoute(router, productController)

	rootRoute := route.NewRootRoute(
		router,
		productRoute,
	)

	rootRoute.Register()

	webContainer := &WebContainer{
		Env:             envConfig,
		ProductDatabase: productDBConfig,
		Repository:      repositoryContainer,
		UseCase:         useCaseContainer,
		Controller:      controllerContainer,
		Route:           rootRoute,
	}

	return webContainer
}
