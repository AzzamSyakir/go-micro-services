package container

import (
	"go-micro-services/src/auth-service/repository"
)

type RepositoryContainer struct {
	Auth *repository.UserRepository
}

func NewRepositoryContainer(
	auth *repository.UserRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		Auth: auth,
	}
	return repositoryContainer
}
