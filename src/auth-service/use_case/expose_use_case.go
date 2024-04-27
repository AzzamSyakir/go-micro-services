package use_case

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-micro-services/src/auth-service/config"
	"go-micro-services/src/auth-service/entity"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	model_response "go-micro-services/src/auth-service/model/response"
	"net/http"
)

type ExposeUseCase struct {
	Env *config.EnvConfig
}

func NewExposeUseCase(envConfig *config.EnvConfig) *ExposeUseCase {
	userUseCase := &ExposeUseCase{
		Env: envConfig,
	}
	return userUseCase
}

// users
func (exposeUseCase *ExposeUseCase) FetchUsers() (result *model_response.Response[[]*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.Host, exposeUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s", address, "users")
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[[]*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseUser
}
func (exposeUseCase *ExposeUseCase) CreateUser(request *model_request.CreateUser) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.Host, exposeUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s", address, "users")
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseUser
}
func (exposeUseCase *ExposeUseCase) DeleteUser(id string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.Host, exposeUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s", address, "users", id)
	newRequest, newRequestErr := http.NewRequest("DELETE", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseUser
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
