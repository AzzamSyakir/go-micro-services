package container

import (
	"go-micro-services/services/Product/repository"
)

type RepositoryContainer struct {
	User *repository.ProductRepository
}

func NewRepositoryContainer(
	user *repository.ProductRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		User: user,
	}
	return repositoryContainer
}
