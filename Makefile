.PHONY: up down test_auth test_product test

up:
	docker-compose up -d --build
down:
	docker-compose down
test_auth:
	@cd ./auth-service/ && go test ./test/
test_product:
	@cd ./product-service/ && go test ./test/
test: test_auth test_product
