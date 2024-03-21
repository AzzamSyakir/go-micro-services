package container

import (
	"fmt"
	"go-micro-services/internal/config"
	http_delivery "go-micro-services/internal/delivery/http"
	"go-micro-services/internal/delivery/http/route"
	"go-micro-services/internal/repository"
	"go-micro-services/internal/use_case"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type WebContainer struct {
	Env        *config.EnvConfig
	Database   *config.DatabaseConfig
	Repository *RepositoryContainer
	UseCase    *UseCaseContainer
	Controller *ControllerContainer
	Route      *route.RootRoute
}

func NewWebContainer() *WebContainer {
	errEnvLoad := godotenv.Load()
	if errEnvLoad != nil {
		panic(fmt.Errorf("error loading .env file: %w", errEnvLoad))
	}

	envConfig := config.NewEnvConfig()
	userDbConfig := config.NewDatabaseConfig(envConfig.UserDB)
	productDbConfig := config.NewDatabaseConfig(envConfig.ProductDB)
	logger := NewLogger(envConfig)
	validate := NewValidator()

	userRepository := repository.NewUserRepository()
	sessionRepository := repository.NewSessionRepository()
	postRepository := repository.NewPostRepository(logger)
	repositoryContainer := NewRepositoryContainer(userRepository, sessionRepository, postRepository)

	userUseCase := use_case.NewUserUseCase(databaseConfig, userRepository, sessionRepository, postRepository)
	authUseCase := use_case.NewAuthUseCase(databaseConfig, userRepository, sessionRepository)
	postUseCase := use_case.NewPostUseCase(databaseConfig, postRepository, logger, validate)
	useCaseContainer := NewUseCaseContainer(userUseCase, authUseCase, postUseCase)

	userController := http.NewUserController(userUseCase)
	postController := http_delivery.NewPostController(postUseCase, logger)
	authController := http_delivery.NewAuthController(authUseCase)
	controllerContainer := NewControllerContainer(userController, authController)

	transactionMiddleware := middleware.NewTransactionMiddleware(databaseConfig)
	authMiddleware := middleware.NewAuthMiddleware(*sessionRepository, databaseConfig)
	rootMiddleware := middleware.NewRootMiddleware(transactionMiddleware, authMiddleware)

	router := mux.NewRouter()
	authRoute := route.NewAuthRoute(router, authController)
	userRoute := route.NewUserRoute(router, userController)
	postRoute := route.NewPostRoute(router, postController)
	protectedRoute := route.NewProtectedRoute(router, userRoute, postRoute)
	unprotectedRoute := route.NewUnprotectedRoute(router, authRoute)
	rootRoute := route.NewRootRoute(
		router,
		rootMiddleware,
		protectedRoute,
		unprotectedRoute,
	)

	rootRoute.Register()

	webContainer := &WebContainer{
		Env:        envConfig,
		Database:   databaseConfig,
		Repository: repositoryContainer,
		UseCase:    useCaseContainer,
		Controller: controllerContainer,
		Middleware: rootMiddleware,
		Route:      rootRoute,
	}

	return webContainer
}
