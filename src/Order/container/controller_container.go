package container

import (
	"go-micro-services/src/Order/delivery/http"
)

type ControllerContainer struct {
	Order *http.OrderController
}

func NewControllerContainer(
	order *http.OrderController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		Order: order,
	}
	return controllerContainer
}
