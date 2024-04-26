package http

import (
	model_request "go-micro-services/src/auth-service/model/request/controller"
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/use_case"
	"net/http"
)

type AuthController struct {
	AuthUseCase *use_case.AuthUseCase
}

func NewAuthController(authUseCase *use_case.AuthUseCase) *AuthController {
	authController := &AuthController{
		AuthUseCase: authUseCase,
	}
	return authController
}
func (authController *AuthController) Login(writer http.ResponseWriter, reader *http.Request) {
	request := model_request.LoginRequest{}
	foundUser, foundUserErr := authController.AuthUseCase.Login(request)
	if foundUserErr == nil {
		response.NewResponse(writer, foundUser)
	}
}
