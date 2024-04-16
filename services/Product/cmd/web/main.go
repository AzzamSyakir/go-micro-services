package main

import (
	"fmt"
	"go-micro-services/services/Product/container"
	"net/http"
)

func main() {
	fmt.Println("Product Services started.")

	webContainer := container.NewWebContainer()

	address := fmt.Sprintf(
		"%s:%s",
		"0.0.0.0",
		webContainer.Env.App.Port,
	)
	listenAndServeErr := http.ListenAndServe(address, webContainer.Route.Router)
	if listenAndServeErr != nil {
		panic(listenAndServeErr)
	}
	fmt.Println("Product Services finished.")
}
