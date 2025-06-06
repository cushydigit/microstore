package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/microstore/product-service/internal/handler"
	"github.com/cushydigit/microstore/product-service/internal/repository"
	"github.com/cushydigit/microstore/product-service/internal/service"
	"github.com/cushydigit/microstore/shared/database"
	"github.com/cushydigit/microstore/shared/middlewares"
	myredis "github.com/cushydigit/microstore/shared/redis"
	"github.com/cushydigit/microstore/shared/zincsearch"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	dsn                = os.Getenv("DSN")
	port               = os.Getenv("PORT")
	redisAddr          = os.Getenv("REDIS_ADDR")
	zincsearchAddr     = os.Getenv("ZINCSEARCH_ADDR")
	zincsearchUsername = os.Getenv("ZINCSEARCH_USERNAME")
	zincsearchPassword = os.Getenv("ZINCSEARCH_PASSWORD")
)

func main() {
	// init redis
	myredis.Init(context.Background(), redisAddr)

	indexer := zincsearch.Init(zincsearchAddr, zincsearchUsername, zincsearchPassword, "products")

	// get dsn
	if dsn == "" {
		log.Panic("DSN not set")
	}
	// connect to db
	db := database.ConnectDB(dsn)

	// TEMP: in-memory product storage
	// repo := repository.NewInMemoryProductRepo()
	repo := repository.NewPostgresProductRepo(db)
	productService := service.NewProductService(repo, indexer)
	productHandler := handler.NewProductHandler(productService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// routes
	r.Route("/product", func(r chi.Router) {
		r.With(middlewares.ValidateCreateProduct).Post("/", productHandler.Create)
		r.Post("/bulk", productHandler.CreateBulk)
		r.Delete("/bulk", productHandler.DeleteAll)
		r.Get("/", productHandler.GetAll)
		r.Get("/search", productHandler.Search)
		r.Get("/{id}", productHandler.GetByID)
		r.Delete("/{id}", productHandler.Delete)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	log.Printf("Starting Product Service on: %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
