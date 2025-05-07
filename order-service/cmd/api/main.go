package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/microstore/order-service/internal/database"
	"github.com/cushydigit/microstore/order-service/internal/handler"
	"github.com/cushydigit/microstore/order-service/internal/repository"
	"github.com/cushydigit/microstore/order-service/internal/service"
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

	repo := repository.NewPostgresOrderRepository(db)
	orderService := service.NewOrderSevice(repo, os.Getenv("PRODUCT_API_URL"))
	orderHandler := handler.NewOrderHandler(orderService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// routes
	r.Route("/order", func(r chi.Router) {
		r.Post("/", orderHandler.Create)
		r.Get("/", orderHandler.GetAll)
		r.Get("/mine", orderHandler.GetByUserID)
		r.Get("/{id}", orderHandler.GetByID)
	})

	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	log.Printf("Starting Order Service on: %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
