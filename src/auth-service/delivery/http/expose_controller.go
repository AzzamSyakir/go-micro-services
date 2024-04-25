package http

import (
	"go-micro-services/src/auth-service/use_case"
	"net/http"
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
func (exposeController *ExposeController) FetchUsers(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) CreateUser(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) DeleteUser(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) UpdateBalance(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) UpdateUser(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) DetailUser(writer http.ResponseWriter, reader *http.Request) {
}

// product

func (exposeController *ExposeController) FetchProducts(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) CreateProduct(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) DeleteProduct(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) UpdateStock(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) UpdateProduct(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) DetailProduct(writer http.ResponseWriter, reader *http.Request) {
}

// category

func (exposeController *ExposeController) FetchCategories(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) CreateCategory(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) DeleteCategory(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) UpdateCategory(writer http.ResponseWriter, reader *http.Request) {
}
func (exposeController *ExposeController) DetailCategory(writer http.ResponseWriter, reader *http.Request) {
}
