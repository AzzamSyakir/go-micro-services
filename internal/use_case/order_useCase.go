package use_case

import (
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"go-micro-services/internal/config"
	"go-micro-services/internal/entity"
	model_request "go-micro-services/internal/model/request/controller"
	model_response "go-micro-services/internal/model/response"
	"go-micro-services/internal/repository"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
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
func (orderUseCase *OrderUseCase) Order(userId string, request model_request.OrderRequest) (result *response.Response[*entity.Order]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
		if err != nil {
			return err
		}
		//  GetoneById User
		userAddress := fmt.Sprintf("%s, %s", orderUseCase.Env.App.Host, orderUseCase.Env.App.Port)
		url := fmt.Sprintf("%s, %s, %s", userAddress, "users", userId)
		requestUser, newRequestErr := http.NewRequest(http.MethodGet, url, nil)
		if requestUser == nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.User]{
				Code:    http.StatusBadRequest,
				Message: "OrderUseCase order is failed, user is not found by id",
				Data:    nil,
			}
			return err
		}
		responseUser, doErr := http.DefaultClient.Do(requestUser)
		if newRequestErr != nil {
			return err
		}
		if doErr != nil {
			return err
		}

		bodyResponseUser := &model_response.Response[*entity.User]{}
		decodeUserErr := json.NewDecoder(responseUser.Body).Decode(bodyResponseUser)
		if decodeUserErr != nil {
			return err
		}
		//  GetoneById Products
		var totalOrderPrice int64
		for i, orderProduct := range request.Products {
			productID := orderProduct.ProductId

			productAdress := fmt.Sprintf("%s, %s", orderUseCase.Env.App.Host, orderUseCase.Env.App.Port)

			productUrl := fmt.Sprintf("%s, %s, %s", productAdress, "products", productID)

			requestProduct, newRequestErr := http.NewRequest(http.MethodGet, productUrl, nil)
			if requestProduct == nil {
				err = begin.Rollback()
				result = &model_response.Response[*entity.User]{
					Code:    http.StatusBadRequest,
					Message: "OrderUseCase order is failed, products is not found by id.",
					Data:    nil,
				}
				return err
			}
			responseProduct, doErr := http.DefaultClient.Do(requestProduct)
			if newRequestErr != nil {
				return err
			}
			if doErr != nil {
				return err
			}

			bodyResponseProduct := &model_response.Response[*entity.Product]{}
			decodeProductErr := json.NewDecoder(responseProduct.Body).Decode(bodyResponseProduct)
			if decodeProductErr != nil {
				return err
			}
			if orderProduct.Qty.Valid && bodyResponseProduct.Data.Stock.Valid {
				if orderProduct.Qty.Int64 > bodyResponseProduct.Data.Stock.Int64 {
					err = begin.Rollback()
					result = &model_response.Response[*entity.User]{
						Code:    http.StatusBadRequest,
						Message: "orderUsecase order is failed, products out off stock",
						Data:    nil,
					}
					return err
				}
			}
			price, _ := strconv.ParseInt(bodyResponseProduct.Data.Price.String, 100, 64)

			totalProductPrice := orderProduct.Qty.Int64 * price
			request.Products[i].TotalPrice.Int64 = totalProductPrice
			totalOrderPrice += totalProductPrice
		}

		totalReturn := request.TotalPaid.Int64 - totalOrderPrice
		TotalPaid := request.TotalPaid.Int64
		//orders
		firstLetter := strings.ToUpper(string(bodyResponseUser.Data.Name.String[0]))
		rand.Seed(time.Now().UnixNano())
		randomDigits := rand.Intn(900) + 100
		receiptCode := fmt.Sprintf("%s%d", firstLetter, randomDigits)
		currentTime := time.Now()
		orderRequest := &entity.Order{
			Id:          null.NewString(uuid.NewString(), true),
			UserId:      bodyResponseUser.Data.Id,
			Name:        bodyResponseUser.Data.Name,
			ReceiptCode: null.NewString(receiptCode, true),
			TotalPrice:  null.NewInt(totalOrderPrice, true),
			TotalPaid:   null.NewInt(TotalPaid, true),
			TotalReturn: null.NewInt(totalReturn, true),
			CreatedAt:   null.NewTime(currentTime, true),
			UpdatedAt:   null.NewTime(currentTime, true),
		}
		orderResult, err := orderUseCase.OrderRepository.Order(begin, orderRequest)
		//	OrderProducts
		var productsInfo []model_request.OrderProducts
		for _, orderProduct := range request.Products {
			orderProductsResult, err := orderUseCase.OrderRepository.OrderProducts(begin, orderProduct)
			if err != nil {
				result = &model_response.Response[*entity.Order]{
					Code:    http.StatusInternalServerError,
					Message: "OrderProductUseCase is failed, order request fail, " + err.Error(),
					Data:    nil,
				}
			}

			productsInfo = orderProduct
		}
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Order]{
			Code:    http.StatusInternalServerError,
			Message: "OrderProductUseCase is failed, order request fail, " + beginErr.Error(),
			Data:    nil,
		}
	}
	return result
}
