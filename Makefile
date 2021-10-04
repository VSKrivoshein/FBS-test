include .env
export $(cat .env | xargs -L 1)

#Только для локальной разработке без запуска приложения без докера
export REDIS_HOST=localhost
export REDIS_PORT=6379

install:
	go mod init github.com/VSKrivoshein/FBS-test && go get ./...

protos:
	protoc -I=internal/app/api/grpc_api/proto --go-grpc_out=internal/app/api/grpc_api/proto --go_out=internal/app/api/grpc_api/proto internal/app/api/grpc_api/proto/*.proto

build:
	 go build ./cmd/apiserver/main.go && ./main

evans:
	evans internal/app/api/proto/fibonacci.proto -p 8081

test:
	docker-compose -f docker-compose.test.yaml up --build

start:
	docker-compose up --build

down:
	docker-compose down

min:
	printenv REDIS_HOST