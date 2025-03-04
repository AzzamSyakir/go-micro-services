package use_case

import (
	"context"
	"fmt"
	"go-micro-services/grpc/pb"
	"go-micro-services/src/order-service/config"
	"go-micro-services/src/order-service/delivery/grpc/client"
	"go-micro-services/src/order-service/repository"
	"math/rand"
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
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve order details: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	orderProducts, err := orderUseCase.OrderRepository.GetOrderProductsByOrderId(begin, id.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve order details: Error retrieving products for the order. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	orderDetails, err := orderUseCase.OrderRepository.DetailOrder(begin, id.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve order details: Database query error. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if orderDetails == nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Order not found: No order exists with the given ID: %s. Rollback status: %v", id.Id, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the database transaction. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	orderDetails.Products = orderProducts.Data

	return &pb.OrderResponse{
		Code:    int64(codes.OK),
		Message: "Successfully retrieved order details.",
		Data:    orderDetails,
	}, nil
}

func (orderUseCase *OrderUseCase) ListOrders(ctx context.Context, _ *pb.Empty) (result *pb.OrderResponseRepeated, err error) {
	begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve orders: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	orderList, err := orderUseCase.OrderRepository.ListOrders(begin)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve orders: Database query error. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if len(orderList.Data) == 0 {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponseRepeated{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("No orders found: There are no orders available in the system. Rollback status: %v", rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	for _, order := range orderList.Data {
		orderProducts, err := orderUseCase.OrderRepository.GetOrderProductsByOrderId(begin, order.Id)
		if err != nil {
			rollbackErr := begin.Rollback()
			return &pb.OrderResponseRepeated{
				Code:    int64(codes.Internal),
				Message: fmt.Sprintf("Failed to retrieve products for order ID %s: Database query error. Error: %v. Rollback status: %v", order.Id, err, rollbackErr),
				Data:    nil,
			}, rollbackErr
		}
		order.Products = orderProducts.Data
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.OrderResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the database transaction. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.OrderResponseRepeated{
		Code:    int64(codes.OK),
		Message: "Successfully retrieved all orders.",
		Data:    orderList.Data,
	}, nil
}

func (orderUseCase *OrderUseCase) Order(ctx context.Context, request *pb.CreateOrderRequest) (result *pb.OrderResponse, err error) {
	begin, err := orderUseCase.DatabaseConfig.OrderDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to create order: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	var totalOrderPrice int64
	for i, product := range request.Products {
		getProduct, err := orderUseCase.productClient.GetProductById(product.ProductId)
		if err != nil {
			rollbackErr := begin.Rollback()
			return &pb.OrderResponse{
				Code:    int64(codes.Internal),
				Message: fmt.Sprintf("Failed to create order: Unable to retrieve product details for Product ID %s. Error: %v. Rollback status: %v", product.ProductId, err, rollbackErr),
				Data:    nil,
			}, rollbackErr
		}
		fmt.Println("productId ", product.ProductId)
		if getProduct.Data == nil {
			rollbackErr := begin.Rollback()
			return &pb.OrderResponse{
				Code:    int64(codes.NotFound),
				Message: fmt.Sprintf("Failed to create order: Product with ID %s not found. %s. Rollback status: %v", product.ProductId, getProduct.Message, rollbackErr),
				Data:    nil,
			}, rollbackErr
		}

		if product.Qty > getProduct.Data.Stock {
			rollbackErr := begin.Rollback()
			return &pb.OrderResponse{
				Code:    int64(codes.FailedPrecondition),
				Message: fmt.Sprintf("Failed to create order: Product with ID %s is out of stock. Requested quantity: %d, Available stock: %d. Rollback status: %v", product.ProductId, product.Qty, getProduct.Data.Stock, rollbackErr),
				Data:    nil,
			}, rollbackErr
		}

		totalProductPrice := product.Qty * getProduct.Data.Price
		request.Products[i].TotalPrice = &totalProductPrice
		totalOrderPrice += totalProductPrice
		finalStock := getProduct.Data.Stock - product.Qty
		orderUseCase.productClient.UpdateProduct(product.ProductId, finalStock)
	}

	getUser, err := orderUseCase.userClient.GetUserById(request.UserId)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to create order: Unable to retrieve user details. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if getUser.Data == nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Failed to create order: User with ID %s not found. %s. Rollback status: %v", request.UserId, getUser.Message, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if request.TotalPaid < totalOrderPrice {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.FailedPrecondition),
			Message: fmt.Sprintf("Failed to create order: Insufficient payment. Total required: %d, Provided: %d. Rollback status: %v", totalOrderPrice, request.TotalPaid, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	finalBalance := getUser.Data.Balance - totalOrderPrice
	orderUseCase.userClient.UpdateUser(request.UserId, finalBalance)

	totalReturn := request.TotalPaid - totalOrderPrice
	firstLetter := strings.ToUpper(string(getUser.Data.Name[0]))
	rand.Seed(time.Now().UnixNano())
	randomDigits := rand.Intn(900) + 100
	receiptCode := fmt.Sprintf("%s%d", firstLetter, randomDigits)

	orderData := &pb.Order{
		Id:          uuid.NewString(),
		UserId:      request.UserId,
		ReceiptCode: receiptCode,
		TotalPrice:  totalOrderPrice,
		TotalPaid:   request.TotalPaid,
		TotalReturn: totalReturn,
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
	}

	order, err := orderUseCase.OrderRepository.Order(begin, orderData)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to create order: Database error while saving order details. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	var productsInfo []*pb.OrderProduct
	for _, orderProduct := range request.Products {
		orderProductData := &pb.OrderProduct{
			Id:         uuid.NewString(),
			OrderId:    order.Data.Id,
			ProductId:  orderProduct.ProductId,
			TotalPrice: totalOrderPrice,
			Qty:        orderProduct.Qty,
			CreatedAt:  timestamppb.New(time.Now()),
			UpdatedAt:  timestamppb.New(time.Now()),
		}
		orderProductSaved, err := orderUseCase.OrderRepository.OrderProducts(begin, orderProductData)
		if err != nil {
			rollbackErr := begin.Rollback()
			return &pb.OrderResponse{
				Code:    int64(codes.Internal),
				Message: fmt.Sprintf("Failed to create order: Error while saving product details for Order ID %s. Error: %v. Rollback status: %v", order.Data.Id, err, rollbackErr),
				Data:    nil,
			}, rollbackErr
		}
		productsInfo = append(productsInfo, orderProductSaved)
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.OrderResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the order transaction. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	result = &pb.OrderResponse{
		Code:    int64(codes.OK),
		Message: "Order created successfully.",
		Data:    order.Data,
	}
	result.Data.Products = productsInfo

	return result, nil
}
