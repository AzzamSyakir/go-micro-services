package container

import (
	"fmt"
	"go-micro-services/src/product-service/config"
	categoryPb "go-micro-services/src/product-service/delivery/grpc/pb/category"
	productPb "go-micro-services/src/product-service/delivery/grpc/pb/product"
	"go-micro-services/src/product-service/repository"
	"go-micro-services/src/product-service/use_case"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type WebContainer struct {
	Env                *config.EnvConfig
	ProductDB          *config.DatabaseConfig
	ProductRepository  *RepositoryContainer
	CategoryRepository *RepositoryContainer
	UseCase            *UseCaseContainer
	Grpc               *grpc.Server
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
	grpcServer := grpc.NewServer()
	productPb.RegisterProductServiceServer(grpcServer, productUseCase)
	categoryPb.RegisterCategoryServiceServer(grpcServer, categoryUseCase)

	webContainer := &WebContainer{
		Env:               envConfig,
		ProductDB:         productDBConfig,
		ProductRepository: repositoryContainer,
		UseCase:           useCaseContainer,
		Grpc:              grpcServer,
	}

	return webContainer
}
