package container

import (
	"go-micro-services/internal/use_case"
)

type UseCaseContainer struct {
	User    *use_case.UserUseCase
	Product *use_case.ProductUseCase
}

func NewUseCaseContainer(
	user *use_case.UserUseCase,
	product *use_case.ProductUseCase,
) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		User:    user,
		Product: product,
	}
	return useCaseContainer
}
