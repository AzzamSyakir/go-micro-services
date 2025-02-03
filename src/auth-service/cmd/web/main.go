package main

import (
	"fmt"
	"go-micro-services/src/auth-service/container"
	"log"
	"net"
	"net/http"
)

func main() {
	fmt.Println("Auth Services started.")

	webContainer := container.NewWebContainer()
	// grpc server
	go func() {
		grpcAddress := fmt.Sprintf(
			"%s:%s",
			"0.0.0.0",
			webContainer.Env.App.AuthGrpcPort,
		)
		netListen, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			log.Fatalf("failed to listen %v", err)
		}
		if err := webContainer.Grpc.Serve(netListen); err != nil {
			log.Fatalf("failed to serve %v", err.Error())
		}
	}()
	// http server
	address := fmt.Sprintf(
		"%s:%s",
		"0.0.0.0",
		webContainer.Env.App.AuthHttpPort,
	)
	listenAndServeErr := http.ListenAndServe(address, webContainer.Middleware.Cors.Handler(webContainer.Route.Router))
	if listenAndServeErr != nil {
		log.Fatalf("failed to serve HTTP: %v", listenAndServeErr)
	}
	fmt.Println("Auth Services finished.")
}
