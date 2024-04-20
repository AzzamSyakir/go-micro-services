package container

import (
	"go-micro-services/src/Product/use_case"
)

type UseCaseContainer struct {
	Product *use_case.ProductUseCase
}

func NewUseCaseContainer(
	product *use_case.ProductUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		Product: product,
	}
	return useCaseContainer
}
