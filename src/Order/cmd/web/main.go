package main

import (
	"fmt"
	"go-micro-services/src/Order/container"
	"net/http"
)

func main() {
	fmt.Println("Order Services started.")

	webContainer := container.NewWebContainer()

	address := fmt.Sprintf(
		"%s:%s",
		"0.0.0.0",
		webContainer.Env.App.OrderPort,
	)
	listenAndServeErr := http.ListenAndServe(address, webContainer.Route.Router)
	if listenAndServeErr != nil {
		panic(listenAndServeErr)
	}
	fmt.Println("Order Services finished.")
}
