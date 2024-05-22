package container

import (
	"fmt"
	"go-micro-services/src/auth-service/delivery/http/route"
	"go-micro-services/src/order-service/config"
	"go-micro-services/src/order-service/delivery/grpc/pb"
	"go-micro-services/src/order-service/repository"
	"go-micro-services/src/order-service/use_case"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type WebContainer struct {
	Env        *config.EnvConfig
	OrderDB    *config.DatabaseConfig
	Repository *RepositoryContainer
	UseCase    *UseCaseContainer
	Route      *route.RootRoute
	Grpc       *grpc.Server
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
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderUseCase)

	webContainer := &WebContainer{
		Env:        envConfig,
		OrderDB:    orderDBConfig,
		Repository: repositoryContainer,
		UseCase:    useCaseContainer,
	}

	return webContainer
}
