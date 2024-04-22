package container

import (
	"fmt"
	"go-micro-services/src/product-service/config"
	httpdelivery "go-micro-services/src/product-service/delivery/http"
	"go-micro-services/src/product-service/delivery/http/route"
	"go-micro-services/src/product-service/repository"
	"go-micro-services/src/product-service/use_case"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
