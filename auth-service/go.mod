module github.com/cushydigit/microstore/auth-service

go 1.24.2

require (
	github.com/cushydigit/microstore/shared v0.1.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cushydigit/microstore/shared => ../shared
