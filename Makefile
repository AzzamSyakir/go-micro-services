start-user:
	go run src/user-service/cmd/web/main.go

start-product:
	go run src/product-service/cmd/web/main.go

start-order:
	go run src/order-service/cmd/web/main.go

start-docker:
	docker compose -f ./docker/docker-compose.yml up -d
stop-docker:
	docker compose -f ./docker/docker-compose.yml down --remove-orphans
start-db:
	docker compose -f ./docker/docker-compose.yml up user-db product-db order-db -d

clean-docker:
	docker system prune && docker volume prune && docker image prune -a -f && docker container prune && docker buildx prune