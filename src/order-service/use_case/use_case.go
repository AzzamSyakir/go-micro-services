package use_case

import (
	"context"
	"fmt"
	"go-micro-services/src/order-service/client"
	"go-micro-services/src/order-service/config"
	"go-micro-services/src/order-service/delivery/grpc/pb"
	"go-micro-services/src/order-service/repository"
	"math/rand"
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
	userClient      *client.UserServiceClient
	productClient   *client.ProductServiceClient
}

func NewOrderUseCase(databaseConfig *config.DatabaseConfig, orderRepository *repository.OrderRepository, envConfig *config.EnvConfig, initUserClient *client.UserServiceClient, initProductClient *client.ProductServiceClient) *OrderUseCase {
	OrderUseCase := &OrderUseCase{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{},
		userClient:                      initUserClient,
		productClient:                   initProductClient,
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
		product, err := orderUseCase.productClient.GetProductById(productId)
		if err != nil {
			rollback := begin.Rollback()
			result = &pb.OrderResponse{
				Code:    int64(codes.Canceled),
				Message: "Order-Service orderUseCase Order is failed, getProduct failed, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}

		if products.Qty > product.Data.Stock {
			rollback := begin.Rollback()
			result = &pb.OrderResponse{
				Code:    int64(codes.Canceled),
				Message: "Order-Service orderUseCase Order is failed, product out of stock.",
				Data:    nil,
			}
			return result, rollback
		}
		totalProductPrice := products.Qty * product.Data.Price
		request.Products[i].TotalPrice = totalProductPrice
		totalOrderPrice += int(totalProductPrice)
		finalStock := product.Data.Stock - products.Qty
		orderUseCase.productClient.UpdateProduct(productId, finalStock)
	}
	//  User
	user, err := orderUseCase.userClient.GetUserById(request.UserId)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderResponse{
			Code:    int64(codes.Canceled),
			Message: "Order-Service orderUseCase Order is failed, GetUser fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	finalBalance := user.Data.Balance - int64(totalOrderPrice)
	orderUseCase.userClient.UpdateUser(request.UserId, finalBalance)
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
	firstLetter := strings.ToUpper(string(user.Data.Name[0]))
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
		requestOrderProducts := &pb.OrderProductRequest{
			ProductId:  orderProducts.ProductId,
			Qty:        orderProducts.Qty,
			TotalPrice: int64(totalOrderPrice),
		}

		orderProduct, _ := orderUseCase.OrderProducts(requestOrderProducts)
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

func (orderUseCase *OrderUseCase) OrderProducts(request *pb.OrderProductRequest) (result *pb.OrderProductResponse, err error) {
	begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.OrderProductResponse{
			Code:    int64(codes.Internal),
			Message: "Order-Service orderUseCase OrderProducts is failed, begin fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	orderProductsData := &pb.OrderProduct{
		Id:         uuid.NewString(),
		OrderId:    request.OrderId,
		ProductId:  request.ProductId,
		TotalPrice: request.TotalPrice,
		Qty:        request.Qty,
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
	commit := begin.Commit()
	result = &pb.OrderProductResponse{
		Data: productsInfo,
	}
	return result, commit
}
