# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY . .
RUN go mod download
RUN go build -o customer-registry-api ./cmd/api
    
# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /src/customer-registry-api .

# TODO: Excesso de privilégios
# Verificar usuário que está executando o container
CMD ["./customer-registry-api"]

