package use_case

import (
	"go-micro-services/src/auth-service/entity"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	model_response "go-micro-services/src/auth-service/model/response"
)

type ExposeUseCase struct {
}

func NewExposeUseCase() *ExposeUseCase {
	userUseCase := &ExposeUseCase{}
	return userUseCase
}

// users
func (exposeUseCase *ExposeUseCase) FetchUsers() (result *model_response.Response[[]*entity.User]) {
	return
}
func (exposeUseCase *ExposeUseCase) CreateUser(request *model_request.CreateUser) (result *model_response.Response[*entity.User]) {
	return
}
func (exposeUseCase *ExposeUseCase) DeleteUser(id string) (result *model_response.Response[*entity.User]) {
	return
}
func (exposeUseCase *ExposeUseCase) UpdateBalance(id string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.User]) {
	return
}
func (exposeUseCase *ExposeUseCase) UpdateUser(id string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.User]) {
	return
}
func (exposeUseCase *ExposeUseCase) DetailUser(id string) (result *model_response.Response[*entity.User]) {
	return
}
func (exposeUseCase *ExposeUseCase) GetOneByEmail(email string) (result *model_response.Response[*entity.User]) {
	return
}

// product

func (exposeUseCase *ExposeUseCase) ListProducts() (result *model_response.Response[[]*entity.Product]) {
	return
}
func (exposeUseCase *ExposeUseCase) CreateProduct(request *model_request.CreateProduct) (result *model_response.Response[*entity.Product]) {
	return
}
func (exposeUseCase *ExposeUseCase) DeleteProduct(id string) (result *model_response.Response[*entity.Product]) {
	return
}
func (exposeUseCase *ExposeUseCase) UpdateStock(id string, request *model_request.ProductPatchOneByIdRequest) (result *model_response.Response[*entity.Product]) {
	return
}
func (exposeUseCase *ExposeUseCase) UpdateProduct(id string, request *model_request.ProductPatchOneByIdRequest) (result *model_response.Response[*entity.Product]) {
	return
}
func (exposeUseCase *ExposeUseCase) DetailProduct(id string) (result *model_response.Response[*entity.Product]) {
	return
}

// category

func (exposeUseCase *ExposeUseCase) ListCategories() (result *model_response.Response[[]*entity.Category]) {
	return
}
func (exposeUseCase *ExposeUseCase) CreateCategory(request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category]) {
	return
}
func (exposeUseCase *ExposeUseCase) DeleteCategory(id string) (result *model_response.Response[*entity.Category]) {
	return
}
func (exposeUseCase *ExposeUseCase) UpdateCategory(id string, request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category]) {
	return
}
func (exposeUseCase *ExposeUseCase) DetailCategory(id string) (result *model_response.Response[*entity.Category]) {
	return
}

// order

func (exposeUseCase *ExposeUseCase) Orders(userId string, request *model_request.OrderRequest) (result *model_response.Response[*model_response.OrderResponse]) {
	return
}
