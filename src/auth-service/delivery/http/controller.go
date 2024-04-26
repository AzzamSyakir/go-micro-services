package http

import (
	"encoding/json"
	"fmt"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/use_case"
	"net/http"
	"strings"
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
	request := &model_request.LoginRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}
	foundUser := authController.AuthUseCase.Login(request)
	response.NewResponse(writer, foundUser)
}
func (authController *AuthController) Logout(writer http.ResponseWriter, reader *http.Request) {
	token := reader.Header.Get("Authorization")
	tokenString := strings.Replace(token, "Bearer ", "", 1)
	fmt.Println(tokenString)

	result := authController.AuthUseCase.Logout(tokenString)
	response.NewResponse(writer, result)
}

func (authController *AuthController) GetNewAccessToken(writer http.ResponseWriter, reader *http.Request) {
	token := reader.Header.Get("Authorization")
	tokenString := strings.Replace(token, "Bearer ", "", 1)

	result := authController.AuthUseCase.GetNewAccessToken(tokenString)
	response.NewResponse(writer, result)
}
