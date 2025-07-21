.PHONY: proto build run docker-build docker-up docker-down

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/pickup.proto

build:
	go build -o bin/pickup-srv ./cmd/server/grpc_server.go

run:
	go run ./cmd/server/grpc_server.go

docker-build:
	docker build -t pickup-srv .

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

grpc-client:
	go build -o bin/grpc_client ./cmd/client/grpc_client.go
	./bin/grpc_client
