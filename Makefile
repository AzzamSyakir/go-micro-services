start-auth:
	clear && go run src/auth-service/cmd/web/main.go

start-user:
	clear && go run src/user-service/cmd/web/main.go

start-product:
	clear && go run src/product-service/cmd/web/main.go

start-order:
	clear && go run src/order-service/cmd/web/main.go

start-docker:
	clear && docker compose -f ./docker/docker-compose.yml up -d 
stop-docker:
	clear && docker compose -f ./docker/docker-compose.yml down --remove-orphans
start-db:
	clear && docker compose -f ./docker/docker-compose.yml up user-db product-db order-db -d

clean-docker:
	clear && docker system prune && docker volume prune && docker image prune -a -f && docker container prune && docker buildx prune