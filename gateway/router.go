package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	authEndpoing    = "http://auth-service:8081"
	productEndpoint = "http://product-service:8082"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	// Middlawares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// specify who is allowed to connect
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // seconds
	}))

	// routes
	// auth service
	r.Post("/login", ProxyHandler(authEndpoing))
	r.Post("/register", ProxyHandler(authEndpoing))

	// product service
	r.Route("/product", func(r chi.Router) {
		// public
		r.Get("/*", ProxyHandler(productEndpoint))

		// private
		r.With(middleware.Logger).Post("/*", ProxyHandler(productEndpoint))
		r.With(middleware.Logger).Delete("/*", ProxyHandler(productEndpoint))
	})

	return r

}
