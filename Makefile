.PHONY: help up down build run test lint coverage migrate-up migrate-down swag tidy

DATABASE_URL ?= postgres://app:app@localhost:5432/customers-registry?sslmode=disable

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Bring up app + postgres
	docker compose up -d

down: ## Tear down stack and remove volumes
	docker compose down -v

build: ## Compile local binary
	go build -o customer-registry-api ./cmd/api

run: ## Run the API locally
	go run ./cmd/api

test: ## Run all tests with race detector
	go test ./...

lint: ## Run golangci-lint
	golangci-lint run

coverage: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-up: ## Apply migrations
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down: ## Roll back last migration	
	migrate -path migrations -database "${DATABASE_URL}" down 1

swag: ## Generate OpenAPI docs
	swag init -g ./cmd/api/main.go -o ./docs

tidy: ## Tidy module dependencies
	go mod tidy
