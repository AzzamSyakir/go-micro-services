package container

import (
	"fmt"
	"go-micro-services/grpc/pb"
	"go-micro-services/src/user-service/config"
	"go-micro-services/src/user-service/delivery/grpc/client"
	"go-micro-services/src/user-service/repository"
	"go-micro-services/src/user-service/use_case"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type WebContainer struct {
	Env        *config.EnvConfig
	UserDB     *config.DatabaseConfig
	Repository *RepositoryContainer
	UseCase    *UseCaseContainer
	Grpc       *grpc.Server
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
	authUrl := fmt.Sprintf(
		"%s:%s",
		envConfig.App.AuthHost,
		envConfig.App.AuthGrpcPort,
	)
	initAuthClient := client.InitAuthServiceClient(authUrl)
	userUseCase := use_case.NewUserUseCase(&initAuthClient, userDBConfig, userRepository)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userUseCase)

	useCaseContainer := NewUseCaseContainer(userUseCase)
	webContainer := &WebContainer{
		Env:        envConfig,
		UserDB:     userDBConfig,
		Repository: repositoryContainer,
		UseCase:    useCaseContainer,
		Grpc:       grpcServer,
	}

	return webContainer
}
