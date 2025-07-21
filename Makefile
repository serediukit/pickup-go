.PHONY: proto build run docker-build docker-up docker-down build-grpc-client build-user-registered-produced

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

build-grpc-client:
	go build -o bin/grpc_client ./cmd/client/grpc_client.go
	./bin/grpc_client

build-user-registered-produced:
	go build -o bin/user_reg_producer ./test/test_producer_user_registered.go
	./bin/user_reg_producer

build-all:
	make docker-up
	make build-user-registered-produced
