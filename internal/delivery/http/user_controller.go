package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	model_request "go-micro-services/internal/model/request/controller"
	"go-micro-services/internal/model/response"
	"go-micro-services/internal/use_case"
	"net/http"
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

func (userController *UserController) PatchOneById(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.UserPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		panic(decodeErr)
	}

	ctx := reader.Context()
	patchedUser, patchedUserErr := userController.UserUseCase.PatchOneByIdFromRequest(ctx, id, request)
	if patchedUserErr == nil {
		response.NewResponse(writer, patchedUser)
	}
}
