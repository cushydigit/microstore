.PHONY: up down reset test_auth test_product test

up:
	docker-compose up -d --build
down:
	docker-compose down
reset: down up
test_auth:
	@cd ./auth-service/ && go test -count=1 ./test/
test_product:
	@cd ./product-service/ && go test -count=1 ./test/
test: test_auth test_product
