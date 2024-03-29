package container

import (
	"go-micro-services/internal/use_case"
)

type UseCaseContainer struct {
	User    *use_case.UserUseCase
	Product *use_case.ProductUseCase
	Order   *use_case.OrderUseCase
}

func NewUseCaseContainer(
	user *use_case.UserUseCase,
	product *use_case.ProductUseCase,
	order *use_case.OrderUseCase,
) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		User:    user,
		Product: product,
		Order:   order,
	}
	return useCaseContainer
}
