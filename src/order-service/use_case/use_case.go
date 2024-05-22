package use_case

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-micro-services/src/order-service/config"
	"go-micro-services/src/order-service/delivery/grpc/pb"
	"go-micro-services/src/order-service/entity"
	model_request "go-micro-services/src/order-service/model/request/controller"
	model_response "go-micro-services/src/order-service/model/response"
	"go-micro-services/src/order-service/repository"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderUseCase struct {
	pb.UnimplementedOrderServiceServer
	DatabaseConfig  *config.DatabaseConfig
	OrderRepository *repository.OrderRepository
	Env             *config.EnvConfig
}

func NewOrderUseCase(databaseConfig *config.DatabaseConfig, orderRepository *repository.OrderRepository, envConfig *config.EnvConfig) *OrderUseCase {
	OrderUseCase := &OrderUseCase{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{},
		DatabaseConfig:                  databaseConfig,
		OrderRepository:                 orderRepository,
		Env:                             envConfig,
	}
	return OrderUseCase
}
func (orderUseCase *OrderUseCase) GetOrderById(ctx context.Context, id *pb.ById) (result *pb.OrderResponse, err error) {
	begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: "order-service DetailOrder is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	orderProductFound, err := orderUseCase.OrderRepository.GetOrderProductsByOrderId(begin, id.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Canceled),
			Message: "order-service DetailOrder is failed, GetOrderProducts fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	orderFound, orderFoundErr := orderUseCase.OrderRepository.DetailOrder(begin, id.Id)
	if orderFoundErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf(": %s", orderFoundErr)
		result = &pb.OrderResponse{
			Code:    int64(codes.Canceled),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	if orderFound == nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Canceled),
			Message: "order-service, DetailOrder is failed, order is not found by id, " + id.Id,
			Data:    nil,
		}

		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.OrderResponse{
		Code:    int64(codes.OK),
		Message: "order-service, DetailOrder is succeed.",
		Data:    orderFound,
	}
	result.Data.Products = orderProductFound.Data
	return result, commit
}

func (orderUseCase *OrderUseCase) ListOrders(context.Context, *pb.Empty) (result *pb.OrderResponseRepeated, err error) {
	begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponseRepeated{
			Code:    int64(codes.Internal),
			Message: "Order-Service orderUseCase ListOrder is failed, begin fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}

	fetchOrder, err := orderUseCase.OrderRepository.ListOrders(begin)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponseRepeated{
			Code:    int64(codes.Canceled),
			Message: "Order-Service orderUseCase ListOrder is failed, query to db  fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	for _, order := range fetchOrder.Data {
		orderProductFound, err := orderUseCase.OrderRepository.GetOrderProductsByOrderId(begin, order.Id)
		if err != nil {
			rollback := begin.Rollback()
			result = &pb.OrderResponseRepeated{
				Code:    int64(codes.Canceled),
				Message: "order-service DetailOrder is failed, GetOrderProducts fail" + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}
		order.Products = orderProductFound.Data
	}
	if fetchOrder.Data == nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponseRepeated{
			Code:    int64(codes.Canceled),
			Message: "orderUseCase ListProduct is failed, data order is empty",
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.OrderResponseRepeated{
		Code:    int64(codes.OK),
		Message: "orderUseCase ListOrder is succeed.",
		Data:    fetchOrder.Data,
	}
	return result, commit
}

func (orderUseCase *OrderUseCase) Order(ctx context.Context, request *pb.Create) (result *pb.OrderResponse, err error) {
	begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: "Order-Service orderUseCase Order is failed, begin fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	//   Products
	var totalOrderPrice int
	for i, products := range request.Products {
		productId := products.ProductId
		product := orderUseCase.GetProduct(productId)
		if product.Data == nil {
			rollback := begin.Rollback()
			result = &pb.OrderResponse{
				Code:    int64(codes.Canceled),
				Message: "Order-Service orderUseCase Order is failed, product  not found.",
				Data:    nil,
			}

			return result, rollback
		}
		if products.Qty > product.Data.Stock.Int64 {
			rollback := begin.Rollback()
			result = &pb.OrderResponse{
				Code:    int64(codes.Canceled),
				Message: "Order-Service orderUseCase Order is failed, product out of stock.",
				Data:    nil,
			}
			return result, rollback
		}
		totalProductPrice := products.Qty * product.Data.Price.Int64
		request.Products[i].TotalPrice = totalProductPrice
		totalOrderPrice += int(totalProductPrice)
		finalStock := product.Data.Stock.Int64 - products.Qty
		orderUseCase.UpdateStock(productId, finalStock)
	}
	//  User
	user := orderUseCase.GetUser(request.UserId)
	if user.Data == nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Canceled),
			Message: "Order-Service orderUseCase Order is failed, user not found.",
			Data:    nil,
		}

		return result, rollback
	}
	finalBalance := user.Data.Balance.Int64 - int64(totalOrderPrice)
	orderUseCase.UpdateBalance(request.UserId, finalBalance)
	//    orders
	if request.TotalPaid < int64(totalOrderPrice) {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Canceled),
			Message: "order-service OrderUseCase is failed, oorder fail,  total paid is not enough, total paid	 required " + string(strconv.FormatInt(int64(totalOrderPrice), 10)),
			Data:    nil,
		}

		return result, rollback
	}
	totalReturn := request.TotalPaid - int64(totalOrderPrice)
	firstLetter := strings.ToUpper(string(user.Data.Name.String[0]))
	rand.Seed(time.Now().UnixNano())
	randomDigits := rand.Intn(900) + 100
	receiptCode := fmt.Sprintf("%s%d", firstLetter, randomDigits)
	orderData := &pb.Order{
		Id:          uuid.NewString(),
		UserId:      request.UserId,
		ReceiptCode: receiptCode,
		TotalPrice:  int64(totalOrderPrice),
		TotalPaid:   request.TotalPaid,
		TotalReturn: totalReturn,
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
		DeletedAt:   timestamppb.New(time.Time{}),
	}

	order, err := orderUseCase.OrderRepository.Order(begin, orderData)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Canceled),
			Message: "order-service OrderUseCase is failed, order  fail,  query to db fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	//    orderProducts
	var productsInfo []*pb.OrderProduct
	for _, orderProducts := range request.Products {
		productId := orderProducts.ProductId
		Qty := orderProducts.Qty

		orderProduct := orderUseCase.OrderProducts(begin, request, productId, Qty, order.Data.Id, totalOrderPrice)
		productsInfoLoop := orderProduct.Data
		productsInfo = append(productsInfo, productsInfoLoop...)
	}

	commit := begin.Commit()

	result = &pb.OrderResponse{
		Code:    int64(codes.OK),
		Message: "orderUseCase success, order is success",
		Data:    order.Data,
	}
	result.Data.Products = productsInfo

	return result, commit
}

func (orderUseCase *OrderUseCase) OrderProducts(begin *sql.Tx, request *model_request.OrderRequest, productId string, Qty int64, orderId string, totalOrderPrice int) (result *pb.OrderProductResponse) {
	orderProductsData := &pb.OrderProduct{
		Id:         uuid.NewString(),
		OrderId:    orderId,
		ProductId:  productId,
		TotalPrice: int64(totalOrderPrice),
		Qty:        Qty,
		CreatedAt:  timestamppb.New(time.Now()),
		UpdatedAt:  timestamppb.New(time.Now()),
		DeletedAt:  timestamppb.New(time.Time{}),
	}
	var productsInfo []*pb.OrderProduct
	orderProduct, orderProductsErr := orderUseCase.OrderRepository.OrderProducts(begin, orderProductsData)
	if orderProductsErr != nil {
		result = &pb.OrderProductResponse{
			Code:    int64(codes.Canceled),
			Message: "orderUseCase fail, order is failed, " + orderProductsErr.Error(),
			Data:    nil,
		}
	}
	productsInfo = append(productsInfo, orderProduct)
	result = &pb.OrderProductResponse{
		Data: productsInfo,
	}
	return result
}

func (orderUseCase *OrderUseCase) GetUser(userId string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.UserHost, orderUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s", address, "users", userId)
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)
	if newRequestErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase failed, GetUser is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase failed, GetUser is failed," + doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    int64(codes.Canceled),
			Message: "orderUseCase fail, GetUser is failed," + decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseUser
}

func (orderUseCase OrderUseCase) UpdateBalance(userId string, balance int64) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.UserHost, orderUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "users", "update-balance", userId)
	payload := map[string]string{"balance": strconv.FormatInt(balance, 10)}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase failed, UpdateBalance user is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase failed, UpdateBalance user is failed," + doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    int64(codes.Canceled),
			Message: "orderUseCase fail, UpdateBalance user is failed," + decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseUser
}

func (orderUseCase *OrderUseCase) GetProduct(productId string) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.ProductHost, orderUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", productId)
	newRequest, newRequestErr := http.NewRequest(http.MethodGet, url, nil)
	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    int64(codes.Canceled),
			Message: "orderUseCase fail, GetProduct is failed, " + newRequestErr.Error(),
			Data:    nil,
		}
	}
	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase failed, GetProduct is failed : " + doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase fail, GetProduct is failed : " + decodeErr.Error(),
			Data:    nil,
		}
		return result
	}
	return bodyResponseProduct
}

func (orderUseCase OrderUseCase) UpdateStock(productId string, stock int64) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", orderUseCase.Env.App.ProductHost, orderUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", productId)
	payload := map[string]string{"stock": strconv.FormatInt(stock, 10)}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase failed, UpdateStock product is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    int64(codes.Canceled),
			Message: "OrderUseCase failed, UpdateStock product is failed," + doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    int64(codes.Canceled),
			Message: "orderUseCase fail, UpdateStock product is failed," + decodeErr.Error(),
			Data:    nil,
		}
	}
	result = nil
	return result
}
