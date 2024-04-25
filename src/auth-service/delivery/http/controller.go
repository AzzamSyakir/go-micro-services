package http

import (
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/use_case"
	"net/http"

	"github.com/gorilla/mux"
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
	vars := mux.Vars(reader)
	id := vars["id"]

	foundUser, foundUserErr := authController.AuthUseCase.GetOneById(id)
	if foundUserErr == nil {
		response.NewResponse(writer, foundUser)
	}
}
