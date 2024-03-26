package container

import (
	"go-micro-services/internal/repository"
)

type RepositoryContainer struct {
	User    *repository.UserRepository
	Product *repository.ProductRepository
	Order   *repository.OrderRepository
}

func NewRepositoryContainer(
	user *repository.UserRepository,
	product *repository.ProductRepository,
	order *repository.OrderRepository,
) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		User:    user,
		Product: product,
		Order:   order,
	}
	return repositoryContainer
}
