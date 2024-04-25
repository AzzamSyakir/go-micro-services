package container

import (
	"go-micro-services/src/auth-service/delivery/http"
)

type ControllerContainer struct {
	Auth *http.AuthController
}

func NewControllerContainer(
	auth *http.AuthController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		Auth: auth,
	}
	return controllerContainer
}
