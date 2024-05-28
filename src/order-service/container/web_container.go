package container

import (
	"fmt"
	"go-micro-services/src/order-service/client"
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
	userUrl := fmt.Sprintf(
		"%s:%s",
		envConfig.App.UserHost,
		envConfig.App.UserPort,
	)
	productUrl := fmt.Sprintf(
		"%s:%s",
		envConfig.App.ProductHost,
		envConfig.App.ProductPort,
	)
	initUserClient := client.InitUserServiceClient(userUrl)
	initProductClient := client.InitProductServiceClient(productUrl)
	orderUseCase := use_case.NewOrderUseCase(orderDBConfig, orderRepository, envConfig, &initUserClient, &initProductClient)

	useCaseContainer := NewUseCaseContainer(orderUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderUseCase)

	webContainer := &WebContainer{
		Env:        envConfig,
		OrderDB:    orderDBConfig,
		Repository: repositoryContainer,
		UseCase:    useCaseContainer,
		Grpc:       grpcServer,
	}

	return webContainer
}
