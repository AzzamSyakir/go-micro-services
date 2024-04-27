package container

import (
	"go-micro-services/src/auth-service/use_case"
)

type UseCaseContainer struct {
	Auth   *use_case.AuthUseCase
	Expose *use_case.ExposeUseCase
}

func NewUseCaseContainer(
	auth *use_case.AuthUseCase,
	expose *use_case.ExposeUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		Auth:   auth,
		Expose: expose,
	}
	return useCaseContainer
}
