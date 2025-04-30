package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/cushydigit/microstore/product-service/internal/database"
	"github.com/cushydigit/microstore/porduct-service/internal/handler"
	"github.com/cushydigit/microstore/porduct-service/internal/repository"
	"github.com/cushydigit/microstore/porduct-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// get dsn
	// dsn := os.Getenv("DSN")
	// if dsn == "" {
	// 	log.Panic("DSN not set")
	// }

	// connect to db
	// db := database.ConnectDB(dsn)

	// TEMP: in-memory product storage
	repo := repository.NewInMemoryProductRepo()
	// repo := repository.PostgresUserRepo(db)
	productService := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(productService)

	r := chi.NewRouter()
	// specify who is allowed to connect
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // seconds
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// routes
	r.Post("/product", productHandler.Create)
	r.Get("/product", productHandler.GetAll)
	r.Get("/product/:id", productHandler.GetByID)
	r.Delete("/product/:id", productHandler.Delete)

	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	log.Printf("Starting Auth Service on: %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
