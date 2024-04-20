package container

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-micro-services/src/User/config"
	httpdelivery "go-micro-services/src/User/delivery/http"
	"go-micro-services/src/User/delivery/http/route"
	"go-micro-services/src/User/repository"
	"go-micro-services/src/User/use_case"
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

	userRepository := repository.NewUserRepository()
	repositoryContainer := NewRepositoryContainer(userRepository)

	userUseCase := use_case.NewUserUseCase(userDBConfig, userRepository)

	useCaseContainer := NewUseCaseContainer(userUseCase)

	userController := httpdelivery.NewUserController(userUseCase)

	controllerContainer := NewControllerContainer(userController)

	router := mux.NewRouter()
	userRoute := route.NewUserRoute(router, userController)

	rootRoute := route.NewRootRoute(
		router,
		userRoute,
	)

	rootRoute.Register()

	webContainer := &WebContainer{
		Env:          envConfig,
		UserDatabase: userDBConfig,
		Repository:   repositoryContainer,
		UseCase:      useCaseContainer,
		Controller:   controllerContainer,
		Route:        rootRoute,
	}

	return webContainer
}
