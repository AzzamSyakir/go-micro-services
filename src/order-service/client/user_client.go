package client

import (
	"context"
	"fmt"
	pb "go-micro-services/src/order-service/delivery/grpc/pb/user"

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
	req := &pb.ByIdUser{
		Id: productId,
	}

	return c.Client.GetUserById(context.Background(), req)
}

func (c *UserServiceClient) UpdateUser(userId string, balance int64) (*pb.UserResponse, error) {
	req := &pb.UpdateUserRequest{
		Id:      userId,
		Balance: &balance,
	}

	return c.Client.UpdateUser(context.Background(), req)
}
