# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.23

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add --no-cache ca-certificates git tzdata && \
    adduser -D -g '' -u 10001 appuser

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/customer-registry-api \
      ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

LABEL org.opencontainers.image.title="customer-registry-api" \
      org.opencontainers.image.source="https://github.com/GuilhermeRossiKirsten/CustomerRegistryAPI" \
      org.opencontainers.image.licenses="MIT"

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /out/customer-registry-api /app/customer-registry-api

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app/customer-registry-api"]
