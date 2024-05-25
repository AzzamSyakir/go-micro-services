package client

import (
	"context"
	"fmt"

	"go-micro-services/src/auth-service/delivery/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryServiceClient struct {
	Client pb.CategoryServiceClient
}

func InitCategoryServiceClient(url string) CategoryServiceClient {
	cc, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := CategoryServiceClient{
		Client: pb.NewCategoryServiceClient(cc),
	}

	return c
}
func (c *CategoryServiceClient) GetCategoryById(productId string) (*pb.CategoryResponse, error) {
	req := &pb.ById{
		Id: productId,
	}

	resp, err := c.Client.GetCategoryById(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}

	return resp, nil
}

func (c *CategoryServiceClient) UpdateCategory(productId string, name string) (*pb.CategoryResponse, error) {
	req := &pb.UpdateCategoryRequest{
		Id:   productId,
		Name: &name,
	}

	resp, err := c.Client.UpdateCategory(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return resp, nil
}
func (c *CategoryServiceClient) CreateCategory(req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {

	resp, err := c.Client.CreateCategory(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return resp, nil
}
func (c *CategoryServiceClient) DeleteCategory(categoryId string) (*pb.CategoryResponse, error) {

	resp, err := c.Client.DeleteCategory(context.Background(), &pb.ById{Id: categoryId})
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	return resp, nil
}
func (c *CategoryServiceClient) ListCategorys(productId string, name string) (*pb.CategoryResponseRepeated, error) {

	resp, err := c.Client.ListCategorys(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return resp, nil
}
