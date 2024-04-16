package container

import (
	"go-micro-services/services/Product/use_case"
)

type UseCaseContainer struct {
	User *use_case.ProductUseCase
}

func NewUseCaseContainer(
	user *use_case.ProductUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		User: user,
	}
	return useCaseContainer
}
