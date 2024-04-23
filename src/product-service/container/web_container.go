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
	Env                *config.EnvConfig
	ProductDatabase    *config.DatabaseConfig
	ProductRepository  *RepositoryContainer
	CategoryRepository *RepositoryContainer
	UseCase            *UseCaseContainer
	Controller         *ControllerContainer
	Route              *route.RootRoute
}

func NewWebContainer() *WebContainer {
	errEnvLoad := godotenv.Load()
	if errEnvLoad != nil {
		panic(fmt.Errorf("error loading .env file: %w", errEnvLoad))
	}

	envConfig := config.NewEnvConfig()
	productDBConfig := config.NewDBConfig(envConfig)

	productRepository := repository.NewProductRepository()
	categoryRepository := repository.NewCategoryRepository()
	repositoryContainer := NewRepositoryContainer(productRepository, categoryRepository)

	productUseCase := use_case.NewProductUseCase(productDBConfig, productRepository)
	categoryUseCase := use_case.NewCategoryUseCase(productDBConfig, categoryRepository)

	useCaseContainer := NewUseCaseContainer(productUseCase, categoryUseCase)

	productController := httpdelivery.NewProductController(productUseCase)
	categoryController := httpdelivery.NewCategoryController(categoryUseCase)

	controllerContainer := NewControllerContainer(productController, categoryController)

	router := mux.NewRouter()
	productRoute := route.NewProductRoute(router, productController)
	categoryRoute := route.NewCategoryRoute(router, categoryController)

	rootRoute := route.NewRootRoute(
		router,
		productRoute,
		categoryRoute,
	)

	rootRoute.Register()

	webContainer := &WebContainer{
		Env:               envConfig,
		ProductDatabase:   productDBConfig,
		ProductRepository: repositoryContainer,
		UseCase:           useCaseContainer,
		Controller:        controllerContainer,
		Route:             rootRoute,
	}

	return webContainer
}
