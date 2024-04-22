package container

import (
	"go-micro-services/src/product-service/delivery/http"
)

type ControllerContainer struct {
	Product *http.ProductController
}

func NewControllerContainer(
	product *http.ProductController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		Product: product,
	}
	return controllerContainer
}
