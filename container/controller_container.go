package container

import (
	"go-micro-services/internal/delivery/http"
)

type ControllerContainer struct {
	User    *http.UserController
	Product *http.ProductController
}

func NewControllerContainer(
	user *http.UserController,
	product *http.ProductController,
) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		User:    user,
		Product: product,
	}
	return controllerContainer
}
