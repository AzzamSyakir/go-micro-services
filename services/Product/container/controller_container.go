package container

import (
	"go-micro-services/services/Product/delivery/http"
)

type ControllerContainer struct {
	User *http.ProductController
}

func NewControllerContainer(
	user *http.ProductController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		User: user,
	}
	return controllerContainer
}
