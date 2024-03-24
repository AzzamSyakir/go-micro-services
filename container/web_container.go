package container

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-micro-services/internal/config"
	http_delivery "go-micro-services/internal/delivery/http"
	"go-micro-services/internal/delivery/http/route"
	"go-micro-services/internal/repository"
	"go-micro-services/internal/use_case"
)

type WebContainer struct {
	Env             *config.EnvConfig
	UserDatabase    *config.DatabaseConfig
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
	userDBConfig := config.NewUserDBConfig(envConfig)
	productDBConfig := config.NewProductDBConfig(envConfig)

	userRepository := repository.NewUserRepository()
	productRepository := repository.NewProductRepository()
	repositoryContainer := NewRepositoryContainer(userRepository, productRepository)

	userUseCase := use_case.NewUserUseCase(userDBConfig, userRepository)
	productUseCase := use_case.NewProductUseCase(productDBConfig, productRepository)
	useCaseContainer := NewUseCaseContainer(userUseCase, productUseCase)

	userController := http_delivery.NewUserController(userUseCase)
	productController := http_delivery.NewProductController(productUseCase)
	controllerContainer := NewControllerContainer(userController, productController)

	router := mux.NewRouter()
	userRoute := route.NewUserRoute(router, userController)
	productRoute := route.NewProductRoute(router, productController)

	rootRoute := route.NewRootRoute(
		router,
		userRoute,
		productRoute,
	)

	rootRoute.Register()

	webContainer := &WebContainer{
		Env:             envConfig,
		UserDatabase:    userDBConfig,
		ProductDatabase: productDBConfig,
		Repository:      repositoryContainer,
		UseCase:         useCaseContainer,
		Controller:      controllerContainer,
		Route:           rootRoute,
	}

	return webContainer
}
