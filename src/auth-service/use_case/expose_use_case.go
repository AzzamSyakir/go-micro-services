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
func (exposeUseCase *ExposeUseCase) ListUsers() (result *model_response.Response[[]*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.UserHost, exposeUseCase.Env.App.UserPort)
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
func (exposeUseCase *ExposeUseCase) CreateUser(request *model_request.Register) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.UserHost, exposeUseCase.Env.App.UserPort)
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
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.UserHost, exposeUseCase.Env.App.UserPort)
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
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.UserHost, exposeUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "users", "update-balance", id)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

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
func (exposeUseCase *ExposeUseCase) UpdateUser(id string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.UserHost, exposeUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s", address, "users", id)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

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
func (exposeUseCase *ExposeUseCase) DetailUser(id string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.UserHost, exposeUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s", address, "users", id)

	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

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
func (exposeUseCase *ExposeUseCase) GetOneByEmail(email string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.UserHost, exposeUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "users", "email", email)

	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

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

// product

func (exposeUseCase *ExposeUseCase) ListProducts() (result *model_response.Response[[]*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s", address, "products")
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[[]*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[[]*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[[]*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[[]*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) CreateProduct(request *model_request.CreateProduct) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s", address, "products")
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) DeleteProduct(id string) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", id)
	newRequest, newRequestErr := http.NewRequest("DELETE", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) UpdateStock(id string, request *model_request.ProductPatchOneByIdRequest) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "products", "update-stock", id)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) UpdateProduct(id string, request *model_request.ProductPatchOneByIdRequest) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", id)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) DetailProduct(id string) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", id)

	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}

// category

func (exposeUseCase *ExposeUseCase) ListCategories() (result *model_response.Response[[]*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s", address, "categories")
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	Category := &model_response.Response[[]*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(Category)
	if decodeErr != nil {
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return Category
}
func (exposeUseCase *ExposeUseCase) CreateCategory(request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s", address, "categories")
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}
func (exposeUseCase *ExposeUseCase) DeleteCategory(id string) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "categories", id)
	newRequest, newRequestErr := http.NewRequest("DELETE", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}
func (exposeUseCase *ExposeUseCase) UpdateCategory(id string, request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "categories", id)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}
func (exposeUseCase *ExposeUseCase) DetailCategory(id string) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "categories", id)

	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}

// order

func (exposeUseCase *ExposeUseCase) Orders(userId string, request *model_request.OrderRequest) (result *model_response.Response[*model_response.OrderResponse]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.OrderHost, exposeUseCase.Env.App.OrderPort)
	url := fmt.Sprintf("%s/%s/%s", address, "orders", userId)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseOrder := &model_response.Response[*model_response.OrderResponse]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseOrder)
	if decodeErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseOrder
}
