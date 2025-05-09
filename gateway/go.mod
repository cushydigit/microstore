module github.com/cushydigit/microstore/gateway

go 1.24.2

require (
	github.com/cushydigit/microstore/shared v0.1.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/cors v1.2.1
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	golang.org/x/crypto v0.38.0 // indirect
)

replace github.com/cushydigit/microstore/shared => ../shared
