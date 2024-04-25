package container

import (
	"go-micro-services/src/auth-service/use_case"
)

type UseCaseContainer struct {
	Auth *use_case.UserUseCase
}

func NewUseCaseContainer(
	auth *use_case.UserUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		Auth: auth,
	}
	return useCaseContainer
}
