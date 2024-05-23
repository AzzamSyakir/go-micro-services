package main

import (
	"fmt"
	"go-micro-services/src/order-service/container"
	"log"
	"net"
)

func main() {
	fmt.Println("Order Services started.")

	webContainer := container.NewWebContainer()

	address := fmt.Sprintf(
		"%s:%s",
		"0.0.0.0",
		webContainer.Env.App.OrderPort,
	)
	netListen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	if err := webContainer.Grpc.Serve(netListen); err != nil {
		log.Fatalf("failed to serve %v", err.Error())
	}
	fmt.Println("Order Services finished.")
}
