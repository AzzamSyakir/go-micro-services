package container

import (
	"go-micro-services/src/product-service/delivery/http"
)

type ControllerContainer struct {
	Product  *http.ProductController
	Category *http.CategoryController
}

func NewControllerContainer(
	product *http.ProductController,
	category *http.CategoryController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		Product:  product,
		Category: category,
	}
	return controllerContainer
}
