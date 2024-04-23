package container

import (
	"go-micro-services/src/product-service/repository"
)

type RepositoryContainer struct {
	Product  *repository.ProductRepository
	Category *repository.CategoryRepository
}

func NewRepositoryContainer(
	product *repository.ProductRepository,
	category *repository.CategoryRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		Product:  product,
		Category: category,
	}
	return repositoryContainer
}
