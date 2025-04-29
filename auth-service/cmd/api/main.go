package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/microstore/auth-service/internal/database"
	"github.com/cushydigit/microstore/auth-service/internal/handler"
	"github.com/cushydigit/microstore/auth-service/internal/repository"
	"github.com/cushydigit/microstore/auth-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// get dsn
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Panic("DSN not set")
	}

	// connect to db
	db := database.ConnectDB(dsn)

	// TEMP: in-memory user storage
	// repo := repository.NewInMemoryUserRepo()
	repo := repository.NewPostgresUserRepo(db)
	authService := service.NewAuthService(repo)
	authHandler := handler.NewAuthHandler(authService)

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
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

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
