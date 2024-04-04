package use_case

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"go-micro-services/internal/config"
	"go-micro-services/internal/entity"
	model_request "go-micro-services/internal/model/request/controller"
	model_response "go-micro-services/internal/model/response"
	"go-micro-services/internal/repository"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type OrderUseCase struct {
	DatabaseConfig  *config.DatabaseConfig
	OrderRepository *repository.OrderRepository
	Env             *config.EnvConfig
}

func NewOrderUseCase(databaseConfig *config.DatabaseConfig, orderRepository *repository.OrderRepository, envConfig *config.EnvConfig) *OrderUseCase {
	OrderUseCase := &OrderUseCase{
		DatabaseConfig:  databaseConfig,
		OrderRepository: orderRepository,
		Env:             envConfig,
	}
	return OrderUseCase
}
func (orderUseCase *OrderUseCase) Order(userId string, request *model_request.OrderRequest) (result *model_response.Response[*model_response.OrderResponse]) {
	begin, beginErr := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
	if beginErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: "orderUseCase fail, order is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}
	//    GetUser
	user := orderUseCase.GetUser(userId)
	if user.Data == nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: user.Message,
			Data:    nil,
		}
		return result
	}
	fmt.Println("user data : ", user.Data)

	//    GetProduct
	var totalOrderPrice int
	for i, products := range request.Products {
		productId := products.ProductId.String
		product := orderUseCase.GetProduct(productId)
		if product.Data == nil {
			result = &model_response.Response[*model_response.OrderResponse]{
				Code:    http.StatusBadRequest,
				Message: product.Message,
				Data:    nil,
			}
			return result
		}
		if products.Qty.Int64 > product.Data.Stock.Int64 {
			result = &model_response.Response[*model_response.OrderResponse]{
				Code:    http.StatusBadRequest,
				Message: "OrderUseCase fail, product out of stock",
				Data:    nil,
			}
			return result
		}
		totalProductPrice := products.Qty.Int64 * product.Data.Price.Int64
		request.Products[i].TotalPrice.Int64 = totalProductPrice
		totalOrderPrice += int(totalProductPrice)
	}
	//    orders
	totalReturn := request.TotalPaid.Int64 - int64(totalOrderPrice)
	firstLetter := strings.ToUpper(string(user.Data.Name.String[0]))
	rand.Seed(time.Now().UnixNano())
	randomDigits := rand.Intn(900) + 100
	receiptCode := fmt.Sprintf("%s%d", firstLetter, randomDigits)
	orderData := &entity.Order{
		Id:          null.NewString(uuid.New().String(), true),
		UserId:      user.Data.Id,
		Name:        user.Data.Name,
		ReceiptCode: null.NewString(receiptCode, true),
		TotalPrice:  null.NewInt(int64(totalOrderPrice), true),
		TotalPaid:   request.TotalPaid,
		TotalReturn: null.NewInt(totalReturn, true),
		CreatedAt:   null.NewTime(time.Now(), true),
		UpdatedAt:   null.NewTime(time.Now(), true),
	}

	order, orderErr := orderUseCase.OrderRepository.Order(begin, orderData)
	if orderErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: "orderUseCase fail, order is failed, " + orderErr.Error(),
			Data:    nil,
		}
		return result
	}
	//    orderProducts
	for _, orderProduct := range request.Products {
		productId := orderProduct.ProductId.String
		orderProductsData := &entity.OrderProducts{
			Id:         null.NewString(uuid.New().String(), true),
			OrderId:    null.NewString(order.Data.Id.String, true),
			ProductId:  null.NewString(productId, true),
			TotalPrice: null.NewInt(int64(totalOrderPrice), true),
			Qty:        null.NewInt(orderProduct.Qty.Int64, true),
			CreatedAt:  null.NewTime(time.Now(), true),
			UpdatedAt:  null.NewTime(time.Now(), true),
		}
		_, orderProductsErr := orderUseCase.OrderRepository.OrderProducts(begin, orderProductsData)
		if orderProductsErr != nil {
			result = &model_response.Response[*model_response.OrderResponse]{
				Code:    http.StatusBadRequest,
				Message: "orderUseCase fail, order is failed, " + orderProductsErr.Error(),
				Data:    nil,
			}
			return result
		}
	}
	order = &model_response.Response[*model_response.OrderResponse]{
		Code:    http.StatusOK,
		Message: "orderUseCase succes, order is success",
		Data:    order.Data,
	}
	return order
}
func (orderUseCase *OrderUseCase) GetUser(userId string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.Port)
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
			Message: "OrderUseCase failed, GetUser is failed," + doErr.Error(),
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
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.Port)
	url := fmt.Sprintf("%s/%s/%s", address, "products", productId)
	fmt.Println(url)
	newRequest, newRequestErr := http.NewRequest(http.MethodGet, url, nil)
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
			Message: "OrderUseCase failed, GetProduct is failed : " + doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase fail, GetProduct is failed : " + decodeErr.Error(),
			Data:    nil,
		}
		return result
	}
	return bodyResponseProduct
}
