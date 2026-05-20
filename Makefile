.PHONY: help up down logs build run test test-race fuzz lint coverage \
        migrate-up migrate-down migrate-create swag tidy \
        docker-build k8s-apply k8s-delete k8s-port-forward

DATABASE_URL    ?= postgres://app:app@localhost:5432/customers-registry?sslmode=disable
IMAGE_NAME      ?= customers-registry-api
IMAGE_TAG       ?= 0.1.0
K8S_MANIFEST    ?= k8s/manifest.yaml

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## ---------- Docker Compose ----------

up: ## Bring up app + postgres (build if needed)
	docker compose up -d --build

down: ## Tear down stack and remove volumes
	docker compose down -v

logs: ## Follow API logs from the compose stack
	docker compose logs -f api

## ---------- Local dev ----------

build: ## Compile local binary
	go build -o customer-registry-api ./cmd/api

run: ## Run the API locally
	go run ./cmd/api

tidy: ## Tidy module dependencies
	go mod tidy

## ---------- Tests ----------

test: ## Run all tests
	go test ./...

test-race: ## Run all tests with race detector
	go test -race ./...

fuzz: ## Run fuzz tests on customer service (30s)
	go test -fuzz=Fuzz -fuzztime=30s ./internal/customer

lint: ## Run golangci-lint
	golangci-lint run

coverage: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## ---------- Migrations ----------

migrate-up: ## Apply all pending migrations
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down: ## Roll back last migration
	migrate -path migrations -database "${DATABASE_URL}" down 1

migrate-create: ## Create a new migration (usage: make migrate-create name=add_phone)
	migrate create -ext sql -dir migrations -seq $(name)

## ---------- Swagger ----------

swag: ## Regenerate OpenAPI docs
	swag init -g ./cmd/api/main.go -o ./docs

## ---------- Docker image ----------

docker-build: ## Build the API container image
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

## ---------- Kubernetes ----------

k8s-apply: ## Apply the k8s manifest
	kubectl apply -f $(K8S_MANIFEST)

k8s-delete: ## Delete resources from the k8s manifest
	kubectl delete -f $(K8S_MANIFEST)

k8s-port-forward: ## Forward localhost:8080 to the API deployment
	kubectl port-forward deploy/customers-registry-api 8080:8080
