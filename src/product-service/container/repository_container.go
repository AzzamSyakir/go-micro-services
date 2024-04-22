package container

import (
	"go-micro-services/src/Product/repository"
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
