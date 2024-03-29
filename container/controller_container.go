package container

import (
	"go-micro-services/internal/delivery/http"
)

type ControllerContainer struct {
	User    *http.UserController
	Product *http.ProductController
	Order   *http.OrderController
}

func NewControllerContainer(
	user *http.UserController,
	product *http.ProductController,
	order *http.OrderController,
) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		User:    user,
		Product: product,
		Order:   order,
	}
	return controllerContainer
}
