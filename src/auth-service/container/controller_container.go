package container

import (
	"go-micro-services/src/auth-service/delivery/http"
)

type ControllerContainer struct {
	Auth *http.UserController
}

func NewControllerContainer(
	auth *http.UserController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		Auth: auth,
	}
	return controllerContainer
}
