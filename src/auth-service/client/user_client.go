package client

import (
	"context"
	"fmt"
	"go-micro-services/src/auth-service/delivery/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	Client pb.UserServiceClient
}

func InitUserServiceClient(url string) UserServiceClient {
	cc, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := UserServiceClient{
		Client: pb.NewUserServiceClient(cc),
	}

	return c
}
func (c *UserServiceClient) GetUserById(productId string) (*pb.UserResponse, error) {
	req := &pb.ById{
		Id: productId,
	}

	resp, err := c.Client.GetUserById(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return resp, nil
}
func (c *UserServiceClient) GetUserByEmail(email string) (*pb.UserResponse, error) {
	req := &pb.ByEmail{
		Email: email,
	}

	return c.Client.GetUserByEmail(context.Background(), req)
}
func (c *UserServiceClient) UpdateUser(userId string, balance int64) (*pb.UserResponse, error) {
	req := &pb.UpdateUserRequest{
		Id:      userId,
		Balance: &balance,
	}

	resp, err := c.Client.UpdateUser(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to Update user: %w", err)
	}

	return resp, nil
}
func (c *UserServiceClient) CreateUser(req *pb.CreateUserRequest) (*pb.UserResponse, error) {

	resp, err := c.Client.CreateUser(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to Create user: %w", err)
	}

	return resp, nil
}
func (c *UserServiceClient) DeleteUser(id string) (*pb.UserResponse, error) {
	resp, err := c.Client.DeleteUser(context.Background(), &pb.ById{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to Delete user: %w", err)
	}

	return resp, nil
}
func (c *UserServiceClient) ListUsers() (*pb.UserResponseRepeated, error) {

	resp, err := c.Client.ListUsers(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to Get List users: %w", err)
	}

	return resp, nil
}
