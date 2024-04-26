package main

import (
	"fmt"
	"go-micro-services/src/auth-service/container"
	"net/http"
)

func main() {
	fmt.Println("Auth Services started.")

	webContainer := container.NewWebContainer()

	address := fmt.Sprintf(
		"%s:%s",
		"0.0.0.0",
		webContainer.Env.App.AuthPort,
	)
	listenAndServeErr := http.ListenAndServe(address, webContainer.Route.Router)
	if listenAndServeErr != nil {
		panic(listenAndServeErr)
	}
	fmt.Println("Auth Services finished.")
}
