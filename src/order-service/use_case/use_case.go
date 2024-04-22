package use_case

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-micro-services/src/order-service/config"
	"go-micro-services/src/order-service/entity"
	model_request "go-micro-services/src/order-service/model/request/controller"
	model_response "go-micro-services/src/order-service/model/response"
	"go-micro-services/src/order-service/repository"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/google/uuid"
	"github.com/guregu/null"
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

	beginErr := crdb.Execute(func() (err error) {
		begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
		if err != nil {
			return err
		}

		//   Products
		var totalOrderPrice int
		for i, products := range request.Products {
			productId := products.ProductId.String
			product := orderUseCase.GetProduct(productId)
			if product.Data == nil {
				err = begin.Rollback()
				result = &model_response.Response[*model_response.OrderResponse]{
					Code:    http.StatusBadRequest,
					Message: "product not found",
					Data:    nil,
				}
				return err
			}
			if products.Qty.Int64 > product.Data.Stock.Int64 {
				err = begin.Rollback()
				result = &model_response.Response[*model_response.OrderResponse]{
					Code:    http.StatusBadRequest,
					Message: "OrderUseCase fail, product out of stock",
					Data:    nil,
				}
				return err
			}
			totalProductPrice := products.Qty.Int64 * product.Data.Price.Int64
			request.Products[i].TotalPrice.Int64 = totalProductPrice
			totalOrderPrice += int(totalProductPrice)
			finalStock := product.Data.Stock.Int64 - products.Qty.Int64
			orderUseCase.UpdateStock(productId, finalStock)
		}
		//  User
		user := orderUseCase.GetUser(userId)
		if user.Data == nil {
			err = begin.Rollback()
			result = &model_response.Response[*model_response.OrderResponse]{
				Code:    http.StatusBadRequest,
				Message: "user not found",
				Data:    nil,
			}
			return err
		}
		finalBalance := user.Data.Balance.Int64 - int64(totalOrderPrice)
		orderUseCase.UpdateBalance(userId, finalBalance)
		//    orders
		if request.TotalPaid.Int64 < int64(totalOrderPrice) {
			err = begin.Rollback()
			result = &model_response.Response[*model_response.OrderResponse]{
				Code:    http.StatusBadRequest,
				Message: "OrderUseCase fail, total payment alone is not enough, total payment required " + string(strconv.FormatInt(int64(totalOrderPrice), 10)),
				Data:    nil,
			}
			return err
		}
		totalReturn := request.TotalPaid.Int64 - int64(totalOrderPrice)
		firstLetter := strings.ToUpper(string(user.Data.Name.String[0]))
		rand.Seed(time.Now().UnixNano())
		randomDigits := rand.Intn(900) + 100
		receiptCode := fmt.Sprintf("%s%d", firstLetter, randomDigits)
		orderData := &entity.Order{
			Id:          null.NewString(uuid.New().String(), true),
			UserId:      user.Data.Id,
			ReceiptCode: null.NewString(receiptCode, true),
			TotalPrice:  null.NewInt(int64(totalOrderPrice), true),
			TotalPaid:   request.TotalPaid,
			TotalReturn: null.NewInt(totalReturn, true),
			CreatedAt:   null.NewTime(time.Now(), true),
			UpdatedAt:   null.NewTime(time.Now(), true),
		}

		order, orderErr := orderUseCase.OrderRepository.Order(begin, orderData)
		if orderErr != nil {
			err = begin.Rollback()
			result = &model_response.Response[*model_response.OrderResponse]{
				Code:    http.StatusBadRequest,
				Message: "orderUseCase fail, order is failed, " + orderErr.Error(),
				Data:    nil,
			}

			return err
		}
		//    orderProducts
		_ = orderUseCase.OrderProducts(begin, request, order.Data.Id.String, totalOrderPrice)

		err = begin.Commit()

		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusOK,
			Message: "orderUseCase success, order is success",
			Data:    order.Data,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusInternalServerError,
			Message: "OrderUseCase order  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}
	return result

}
func (orderUseCase *OrderUseCase) GetUser(userId string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s", address, "User", userId)
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

func (orderUseCase OrderUseCase) UpdateBalance(userId string, balance int64) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "User", "update-balance", userId)
	payload := map[string]string{"balance": strconv.FormatInt(balance, 10)}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, UpdateBalance user is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, UpdateBalance user is failed," + doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "orderUseCase fail, UpdateBalance user is failed," + decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseUser
}

func (orderUseCase *OrderUseCase) GetProduct(productId string) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", productId)
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

func (orderUseCase *OrderUseCase) OrderProducts(begin *sql.Tx, request *model_request.OrderRequest, orderId string, totalOrderPrice int) (result *model_response.Response[*model_response.OrderResponse]) {
	for _, orderProduct := range request.Products {
		productId := orderProduct.ProductId.String
		orderProductsData := &entity.OrderProducts{
			Id:         null.NewString(uuid.New().String(), true),
			OrderId:    null.NewString(orderId, true),
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
		}
	}
	result = nil
	return result
}

func (orderUseCase OrderUseCase) UpdateStock(productId string, stock int64) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.Host, orderUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "products", "update-stock", productId)
	payload := map[string]string{"stock": strconv.FormatInt(stock, 10)}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, UpdateStock product is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, UpdateStock product is failed," + doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: "orderUseCase fail, UpdateStock product is failed," + decodeErr.Error(),
			Data:    nil,
		}
	}
	result = nil
	return result

}
