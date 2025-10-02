# Env file
ENV_FILE := .env

# Docker compose paths
DOCKER_COMPOSE_DEV := ./docker/docker-compose.dev.yml
DOCKER_COMPOSE_PROD := ./docker/docker-compose.prod.yml

# services command
start-auth:
	clear && go run src/auth-service/cmd/web/main.go

start-user:
	clear && go run src/user-service/cmd/web/main.go

start-product:
	clear && go run src/product-service/cmd/web/main.go

start-order:
	clear && go run src/order-service/cmd/web/main.go

# docker command
up-dev:
	clear && docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_DEV) up -d 

down-dev:
	clear && docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_DEV) down -v --remove-orphans

up-prod:
	clear && docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_PROD) up -d 

down-prod:
	clear && docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_PROD) down -v --remove-orphans

clean-docker:
	clear && docker system prune -f && docker volume prune -f && docker image prune -a -f && docker container prune -f && docker buildx prune -f

start-test:
	clear && go test -v -count=1 ./src/auth-service/test

generate-proto:
	clear && protoc --proto_path=grpc/proto grpc/proto/*.proto --go_out=grpc --go-grpc_out=grpc
