package client

import (
	"context"
	"fmt"

	pb "go-micro-services/src/order-service/delivery/grpc/pb/product"

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
	req := &pb.ByIdProduct{
		Id: productId,
	}

	return c.Client.GetProductById(context.Background(), req)
}

func (c *ProductServiceClient) UpdateProduct(productId string, stock int64) (*pb.ProductResponse, error) {
	req := &pb.UpdateProductRequest{
		Id:    productId,
		Stock: &stock,
	}

	return c.Client.UpdateProduct(context.Background(), req)
}
