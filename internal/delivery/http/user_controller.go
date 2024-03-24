package http

import (
	"github.com/gorilla/mux"
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

	result, resultErr := userController.UserUseCase.GetOneById(id)
	if resultErr != nil {
		response.NewResponse(writer, result)
	}
	response.NewResponse(writer, result)
}
