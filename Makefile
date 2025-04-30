.PHONY: up down test_auth test

up:
	docker-compose up -d --build
down:
	docker-compose down
test_auth:
	@cd ./auth-service/ && go test ./test/
test: test_auth
