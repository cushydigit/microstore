package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/microstore/product-service/internal/database"
	"github.com/cushydigit/microstore/product-service/internal/handler"
	"github.com/cushydigit/microstore/product-service/internal/repository"
	"github.com/cushydigit/microstore/product-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// get dsn
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Panic("DSN not set")
	}

	// connect to db
	db := database.ConnectDB(dsn)

	// TEMP: in-memory product storage
	// repo := repository.NewInMemoryProductRepo()
	repo := repository.NewPostgresProductRepo(db)
	productService := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(productService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// routes
	r.Route("/product", func(r chi.Router) {
		r.Post("/", productHandler.Create)
		r.Get("/", productHandler.GetAll)
		r.Get("/{id}", productHandler.GetByID)
		r.Delete("/{id}", productHandler.Delete)
	})

	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	log.Printf("Starting Product Service on: %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
