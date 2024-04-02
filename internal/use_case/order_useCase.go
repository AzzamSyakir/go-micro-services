package use_case

import (
	"encoding/json"
	"fmt"
	"go-micro-services/internal/config"
	"go-micro-services/internal/entity"
	model_request "go-micro-services/internal/model/request/controller"
	model_response "go-micro-services/internal/model/response"
	"go-micro-services/internal/repository"
	"net/http"
)

type OrderUseCase struct {
	DatabaseConfig  *config.DatabaseConfig
	OrderRepository *repository.OrderRepository
	Env             *config.EnvConfig
}

func NewOrderUseCase(databaseConfig *config.DatabaseConfig, orderRepository *repository.OrderRepository) *OrderUseCase {
	OrderUseCase := &OrderUseCase{
		DatabaseConfig:  databaseConfig,
		OrderRepository: orderRepository,
	}
	return OrderUseCase
}
func (orderUseCase *OrderUseCase) Order(userId string, request *model_request.OrderRequest) (result *model_response.Response[*model_response.OrderResponse], err error) {
	//	GetUser
	user := orderUseCase.GetUser(userId)
	//	GetProduct
	for i, products := range request.Products {
		productId := products.ProductId.String
		product := orderUseCase.GetProduct(productId)
		if products.Qty.Int64 > product.Data.Stock.Int64 {
			result = &model_response.Response[*model_response.OrderResponse]{
				Code:    http.StatusBadRequest,
				Message: "OrderUseCase fail, product out of stock",
				Data:    nil,
			}
		}
		price := products.Qty.Int64 * product.Data.Price.Int64
		request.Products[i].TotalPrice.Int64 = price
	}
}
func (orderUseCase *OrderUseCase) GetUser(userId string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.Port)
	url := fmt.Sprintf("%s/%s/%s", address, "users", userId)
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)
	if newRequestErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, GetUser is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, GetUser is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "orderUseCase fail, GetUser is failed," + decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseUser
}

func (orderUseCase *OrderUseCase) GetProduct(productId string) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.Port)
	url := fmt.Sprintf("%s/%s/%s", address, "users", productId)
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)
	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: "orderUseCase fail, GetProduct is failed, " + newRequestErr.Error(),
			Data:    nil,
		}
	}
	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, GetProduct is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase fail, GetProduct is failed, " + decodeErr.Error(),
			Data:    nil,
		}
		return result
	}
	return bodyResponseProduct
}
