package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ErrorJSON(w, errors.New("Authorization header missing"), http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ErrorJSON(w, errors.New("Invalid token format"), http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			ErrorJSON(w, errors.New("Invalid or expired token"), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			ErrorJSON(w, errors.New("Invalid token claims"), http.StatusUnauthorized)
			return
		}

		userID := int(claims["user_id"].(float64))

		// Inject int o header fo downstream services
		r.Header.Set("X-User-ID", string(rune(userID)))

		// Inject into context in the case of future use
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDVal := r.Context().Value("user_id")
		userID, ok := userIDVal.(int)
		if !ok || userID != 1 {
			ErrorJSON(w, errors.New("Require admin privilate"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
