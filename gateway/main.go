package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	myredis "github.com/cushydigit/microstore/shared/redis"
)

var (
	authEndpoint    = os.Getenv("AUTH_API_URL")
	productEndpoint = os.Getenv("PRODUCT_API_URL")
	orderEndpoint   = os.Getenv("ORDER_API_URL")
	redisAddr       = os.Getenv("REDIS_ADDR")
	port            = os.Getenv("PORT")
)

func main() {

	myredis.Init(redisAddr)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: Routes(),
	}

	log.Printf("API Gateway listening on: %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
