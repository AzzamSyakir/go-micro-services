.PHONY: build test docker

build-user:
	CGO_ENABLED=0 go build -a -o User user-service/cmd/server/main.go

build-payment: 
	CGO_ENABLED=0 go build -a -o payment payment-service/cmd/server/main.go

test: 
	go test -v ./... -cover -vet -all
