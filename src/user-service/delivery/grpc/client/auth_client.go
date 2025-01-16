package client

import (
	"context"
	"fmt"
	"go-micro-services/grpc/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Client pb.AuthServiceClient
}

func InitAuthServiceClient(url string) AuthServiceClient {
	cc, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := AuthServiceClient{
		Client: pb.NewAuthServiceClient(cc),
	}
	fmt.Println("init auth grpc service", url)
	return c
}
func (c *AuthServiceClient) LogoutWithUserId(req *pb.ByUserId) (*pb.Empty, error) {
	resp, err := c.Client.LogoutWithUserId(context.Background(), req)
	if err != nil {
		log.Fatal("failed to LogoutWithUserId: %w", err)
	}
	return resp, nil
}
