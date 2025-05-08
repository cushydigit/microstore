# product-service/Dockerfile => require to build from the root
FROM docker.arvancloud.ir/golang:1.24-alpine AS builder

WORKDIR /app

COPY shared ./shared
COPY product-service ./product-service

# Download dependencies for shared module
WORKDIR /app/shared
RUN go mod download

# Download and build service
WORKDIR /app/product-service
RUN go clean -modcache
RUN go mod tidy
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -o /bin/product-service ./cmd/api/

# --- Final Image ---
FROM docker.arvancloud.ir/alpine:latest

WORKDIR /app

COPY --from=builder /bin/product-service .

CMD ["/app/product-service"]

