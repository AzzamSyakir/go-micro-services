package container

import (
	"go-micro-services/internal/repository"
)

type RepositoryContainer struct {
	User    *repository.UserRepository
	Product *repository.ProductRepository
}

func NewRepositoryContainer(
	user *repository.UserRepository,
	Product *repository.ProductRepository,
) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		User:    user,
		Product: Product,
	}
	return repositoryContainer
}
