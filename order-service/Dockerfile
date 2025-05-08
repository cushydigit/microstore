# order-service/Dockerfile => require to build from the root
FROM docker.arvancloud.ir/golang:1.24-alpine AS builder

WORKDIR /app

COPY shared ./shared
COPY order-service ./order-service

# Download dependencies for shared module
WORKDIR /app/shared
RUN go mod download

# Download and build service
WORKDIR /app/order-service
RUN go clean -modcache
RUN go mod tidy
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -o /bin/order-service ./cmd/api/

# --- Final Image ---
FROM docker.arvancloud.ir/alpine:latest

WORKDIR /app

COPY --from=builder /bin/order-service .

CMD ["/app/order-service"]

