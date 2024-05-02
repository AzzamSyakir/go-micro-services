package http

import (
	"encoding/json"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/use_case"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ExposeController struct {
	ExposeUseCase *use_case.ExposeUseCase
}

func NewExposeController(exposeUseCase *use_case.ExposeUseCase) *ExposeController {
	exposeController := &ExposeController{
		ExposeUseCase: exposeUseCase,
	}
	return exposeController
}

// users

func (exposeController *ExposeController) FetchUser(writer http.ResponseWriter, reader *http.Request) {
	fetchUser := exposeController.ExposeUseCase.ListUsers()
	response.NewResponse(writer, fetchUser)
}
func (exposeController *ExposeController) Register(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.RegisterRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := exposeController.ExposeUseCase.CreateUser(request)

	response.NewResponse(writer, result)
}
func (exposeController *ExposeController) DeleteUser(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	result := exposeController.ExposeUseCase.DeleteUser(id)

	response.NewResponse(writer, result)
}
func (exposeController *ExposeController) UpdateBalance(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.UserPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := exposeController.ExposeUseCase.UpdateBalance(id, request)

	response.NewResponse(writer, result)
}
func (exposeController *ExposeController) UpdateUser(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.UserPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := exposeController.ExposeUseCase.UpdateUser(id, request)

	response.NewResponse(writer, result)
}
func (expoaseController *ExposeController) DetailUser(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	foundUser := expoaseController.ExposeUseCase.DetailUser(id)
	response.NewResponse(writer, foundUser)
}
func (exposeController *ExposeController) GetUserByEmail(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	email := vars["email"]

	foundUser := exposeController.ExposeUseCase.GetOneByEmail(email)
	response.NewResponse(writer, foundUser)
}

// product

func (exposeController *ExposeController) ListProducts(writer http.ResponseWriter, reader *http.Request) {
	product := exposeController.ExposeUseCase.ListProducts()
	response.NewResponse(writer, product)
}
func (exposeController *ExposeController) CreateProduct(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.CreateProduct{}

	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := exposeController.ExposeUseCase.CreateProduct(request)

	response.NewResponse(writer, result)
}

func (exposeController *ExposeController) DeleteProduct(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	result := exposeController.ExposeUseCase.DeleteProduct(id)

	response.NewResponse(writer, result)
}
func (exposeController *ExposeController) UpdateStock(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.ProductPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		panic(decodeErr)
	}
	result := exposeController.ExposeUseCase.UpdateStock(id, request)

	response.NewResponse(writer, result)
}
func (exposeController *ExposeController) UpdateProduct(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.ProductPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		panic(decodeErr)
	}
	result := exposeController.ExposeUseCase.UpdateProduct(id, request)

	response.NewResponse(writer, result)
}
func (exposeController *ExposeController) DetailProduct(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	foundProduct := exposeController.ExposeUseCase.DetailProduct(id)
	response.NewResponse(writer, foundProduct)
}

// category

func (exposeController *ExposeController) CreateCategory(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.CategoryRequest{}

	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, "Failed to decode request body: "+decodeErr.Error(), http.StatusBadRequest)
		return
	}

	result := exposeController.ExposeUseCase.CreateCategory(request)

	response.NewResponse(writer, result)
}

func (exposeController *ExposeController) ListCategories(writer http.ResponseWriter, reader *http.Request) {
	Category := exposeController.ExposeUseCase.ListCategories()
	response.NewResponse(writer, Category)
}

func (exposeController *ExposeController) DeleteCategory(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	result := exposeController.ExposeUseCase.DeleteCategory(id)

	response.NewResponse(writer, result)
}

func (exposeController *ExposeController) UpdateCategory(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.CategoryRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		panic(decodeErr)
	}
	result := exposeController.ExposeUseCase.UpdateCategory(id, request)

	response.NewResponse(writer, result)
}
func (exposeController *ExposeController) DetailCategory(writer http.ResponseWriter, reader *http.Request) {
}

// order

func (exposeController *ExposeController) Orders(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.OrderRequest{}
	token := reader.Header.Get("Authorization")
	tokenString := strings.Replace(token, "Bearer ", "", 1)
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, "Failed to decode request body: "+decodeErr.Error(), http.StatusBadRequest)
		return
	}
	result := exposeController.ExposeUseCase.Orders(tokenString, request)
	response.NewResponse(writer, result)
}
