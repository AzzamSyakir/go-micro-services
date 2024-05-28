package client

import (
	"context"
	"fmt"

	"go-micro-services/src/auth-service/delivery/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderServiceClient struct {
	Client pb.OrderServiceClient
}

func InitOrderServiceClient(url string) OrderServiceClient {
	cc, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := OrderServiceClient{
		Client: pb.NewOrderServiceClient(cc),
	}

	return c
}
func (c *OrderServiceClient) GetOrderById(productId string) (*pb.OrderResponse, error) {
	req := &pb.ById{
		Id: productId,
	}

	resp, err := c.Client.GetOrderById(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	return resp, nil
}
func (c *OrderServiceClient) Order(req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {

	resp, err := c.Client.Order(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to order: %w", err)
	}

	return resp, nil
}
func (c *OrderServiceClient) ListOrders() (*pb.OrderResponseRepeated, error) {

	resp, err := c.Client.ListOrders(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to get list orders: %w", err)
	}

	return resp, nil
}
