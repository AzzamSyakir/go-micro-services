package container

import (
	"go-micro-services/services/User/repository"
)

type RepositoryContainer struct {
	User *repository.UserRepository
}

func NewRepositoryContainer(
	user *repository.UserRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		User: user,
	}
	return repositoryContainer
}
