package container

import (
	"go-micro-services/src/product-service/repository"
)

type RepositoryContainer struct {
	Product *repository.ProductRepository
}

func NewRepositoryContainer(
	product *repository.ProductRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		Product: product,
	}
	return repositoryContainer
}
