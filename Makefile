.PHONY: help up down build run test lint coverage migrate-up migrate-down swag tidy

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Bring up app + postgres

down: ## Tear down stack and remove volumes

build: ## Compile local binary

run: ## Run the API locally

test: ## Run all tests with race detector

lint: ## Run golangci-lint

coverage: ## Generate HTML coverage report

migrate-up: ## Apply migrations

migrate-down: ## Roll back last migration	

swag: ## Generate OpenAPI docs

tidy: ## Tidy module dependencies
