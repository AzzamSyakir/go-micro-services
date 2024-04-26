package http

import (
	"encoding/json"
	model_request "go-micro-services/src/user-service/model/request/controller"
	"go-micro-services/src/user-service/model/response"
	"go-micro-services/src/user-service/use_case"
	"net/http"

	"github.com/gorilla/mux"
)

type UserController struct {
	UserUseCase *use_case.UserUseCase
}

func NewUserController(userUseCase *use_case.UserUseCase) *UserController {
	userController := &UserController{
		UserUseCase: userUseCase,
	}
	return userController
}
func (userController *UserController) GetOneById(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	foundUser, foundUserErr := userController.UserUseCase.GetOneById(id)
	if foundUserErr == nil {
		response.NewResponse(writer, foundUser)
	}
}
func (userController *UserController) GetOneByEmail(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	email := vars["email"]

	foundUser, foundUserErr := userController.UserUseCase.GetOneByEmail(email)
	if foundUserErr == nil {
		response.NewResponse(writer, foundUser)
	}
}
func (userController *UserController) FetchUser(writer http.ResponseWriter, reader *http.Request) {
	fetchUser, fetchUserErr := userController.UserUseCase.FetchUser()
	if fetchUserErr == nil {
		response.NewResponse(writer, fetchUser)
	}
}

func (userController *UserController) UpdateBalance(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.UserPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := userController.UserUseCase.UpdateBalance(id, request)

	response.NewResponse(writer, result)
}
func (userController *UserController) UpdateUser(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.UserPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := userController.UserUseCase.UpdateUser(id, request)

	response.NewResponse(writer, result)
}

func (userController *UserController) CreateUser(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.CreateUser{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := userController.UserUseCase.CreateUser(request)

	response.NewResponse(writer, result)
}
func (userController *UserController) DeleteUser(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	result := userController.UserUseCase.DeleteUser(id)

	response.NewResponse(writer, result)
}
