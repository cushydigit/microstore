module github.com/cushydigit/microstore/order-service

go 1.24.2

require (
	github.com/cushydigit/microstore/shared v0.1.0
	github.com/go-chi/chi/v5 v5.2.1
)

require github.com/lib/pq v1.10.9 // indirect

replace github.com/cushydigit/microstore/shared => ../shared
