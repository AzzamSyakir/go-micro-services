package http

import (
	"encoding/json"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/use_case"
	"net/http"
	"strings"
)

type AuthController struct {
	AuthUseCase   *use_case.AuthUseCase
	ExposeUseCase *use_case.ExposeUseCase
}

func NewAuthController(authUseCase *use_case.AuthUseCase, exposeUseCase *use_case.ExposeUseCase) *AuthController {
	authController := &AuthController{
		AuthUseCase:   authUseCase,
		ExposeUseCase: exposeUseCase,
	}
	return authController
}
func (authController *AuthController) Register(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.RegisterRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := authController.ExposeUseCase.CreateUser(request)

	response.NewResponse(writer, result)
}
func (authController *AuthController) Login(writer http.ResponseWriter, reader *http.Request) {
	request := &model_request.LoginRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}
	foundUser, _ := authController.AuthUseCase.Login(request)
	response.NewResponse(writer, foundUser)
}
func (authController *AuthController) Logout(writer http.ResponseWriter, reader *http.Request) {
	token := reader.Header.Get("Authorization")
	tokenString := strings.Replace(token, "Bearer ", "", 1)

	result := authController.AuthUseCase.Logout(tokenString)
	response.NewResponse(writer, result)
}

func (authController *AuthController) GetNewAccessToken(writer http.ResponseWriter, reader *http.Request) {
	token := reader.Header.Get("Authorization")
	tokenString := strings.Replace(token, "Bearer ", "", 1)

	result := authController.AuthUseCase.GetNewAccessToken(tokenString)
	response.NewResponse(writer, result)
}
