package container

import (
	"go-micro-services/src/product-service/use_case"
)

type UseCaseContainer struct {
	product-service *use_case.ProductUseCase
}

func NewUseCaseContainer(
	product *use_case.ProductUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		product-service: product,
	}
	return useCaseContainer
}