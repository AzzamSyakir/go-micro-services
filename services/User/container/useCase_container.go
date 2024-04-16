package container

import (
	"go-micro-services/services/User/use_case"
)

type UseCaseContainer struct {
	User *use_case.UserUseCase
}

func NewUseCaseContainer(
	user *use_case.UserUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		User: user,
	}
	return useCaseContainer
}
