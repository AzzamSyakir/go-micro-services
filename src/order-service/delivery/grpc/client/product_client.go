package client

import (
	"context"
	"fmt"

	"go-micro-services/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductServiceClient struct {
	Client pb.ProductServiceClient
}

func InitProductServiceClient(url string) ProductServiceClient {
	cc, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := ProductServiceClient{
		Client: pb.NewProductServiceClient(cc),
	}

	return c
}
func (c *ProductServiceClient) GetProductById(productId string) (*pb.ProductResponse, error) {
	req := &pb.ById{
		Id: productId,
	}

	resp, err := c.Client.GetProductById(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return resp, nil
}

func (c *ProductServiceClient) UpdateProduct(productId string, stock int64) (*pb.ProductResponse, error) {
	req := &pb.UpdateProductRequest{
		Id:    productId,
		Stock: &stock,
	}

	resp, err := c.Client.UpdateProduct(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return resp, nil
}
