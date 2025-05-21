package main

import (
	"github.com/cushydigit/microstore/shared/middlewares"
	"github.com/cushydigit/microstore/shared/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.RateLimiter)

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
	r.Post("/login", utils.ProxyHandler(authEndpoint))
	r.Post("/register", utils.ProxyHandler(authEndpoint))

	// product service
	r.Route("/product", func(r chi.Router) {
		// public
		r.Get("/*", utils.ProxyHandler(productEndpoint))

		// private admin only
		r.With(middlewares.RequireAuth, middlewares.RequireAdmin).Post("/*", utils.ProxyHandler(productEndpoint))
		r.With(middlewares.RequireAuth, middlewares.RequireAdmin).Delete("/*", utils.ProxyHandler(productEndpoint))
	})

	// order service
	r.Route("/order", func(r chi.Router) {
		// private admin only
		r.With(middlewares.RequireAuth, middlewares.RequireAdmin).Get("/", utils.ProxyHandler(orderEndpoint))
		// private user
		r.With(middlewares.RequireAuth).Post("/", utils.ProxyHandler(orderEndpoint))
		r.With(middlewares.RequireAuth).Get("/mine", utils.ProxyHandler(orderEndpoint))
		r.With(middlewares.RequireAuth).Get("/{id}", utils.ProxyHandler(orderEndpoint))
	})

	return r

}
