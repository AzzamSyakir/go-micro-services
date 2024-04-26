package container

import (
	"go-micro-services/src/auth-service/use_case"
)

type UseCaseContainer struct {
	Auth *use_case.AuthUseCase
}

func NewUseCaseContainer(
	auth *use_case.AuthUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		Auth: auth,
	}
	return useCaseContainer
}
