package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	authEndpoint    = os.Getenv("AUTH_API_URL")
	productEndpoint = os.Getenv("PRODUCT_API_URL")
	orderEndpoint   = os.Getenv("ORDER_API_URL")
)

func Routes() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// specify who is allowed to connect
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5134",  // Dev frontend
			"https://microstore.com", // Prod frontend
		},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{
			"Authorization",
			"X-User-ID",
		},
		AllowCredentials: true,
		MaxAge:           300, // seconds
	}))

	// routes
	// auth service
	r.Post("/login", ProxyHandler(authEndpoint))
	r.Post("/register", ProxyHandler(authEndpoint))

	// product service
	r.Route("/product", func(r chi.Router) {
		// public
		r.Get("/*", ProxyHandler(productEndpoint))

		// private
		r.With(AuthMiddleware, AdminMiddleware).Post("/*", ProxyHandler(productEndpoint))
		r.With(AuthMiddleware, AdminMiddleware).Delete("/*", ProxyHandler(productEndpoint))
	})

	// order service
	r.Route("/order", func(r chi.Router) {
		// private admin route
		r.With(AuthMiddleware, AdminMiddleware).Get("/", ProxyHandler(orderEndpoint))
		// private user route
		r.With(AuthMiddleware).Post("/", ProxyHandler(orderEndpoint))
		r.With(AuthMiddleware).Get("/mine", ProxyHandler(orderEndpoint))
		r.With(AuthMiddleware).Get("/{id}", ProxyHandler(orderEndpoint))
	})

	return r

}
