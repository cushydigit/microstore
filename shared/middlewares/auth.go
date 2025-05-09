package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cushydigit/microstore/shared/helpers"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.ErrorJSON(w, errors.New("Authorization header missing"), http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			helpers.ErrorJSON(w, errors.New("Invalid token format"), http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			helpers.ErrorJSON(w, errors.New("Invalid or expired token"), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			helpers.ErrorJSON(w, errors.New("Invalid token claims"), http.StatusUnauthorized)
			return
		}

		userID := int(claims["user_id"].(float64))

		// Inject int o header fo downstream services
		r.Header.Set("X-User-ID", strconv.Itoa(userID))

		// Inject into context in the case of future use
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDVal := r.Context().Value("user_id")
		userID, ok := userIDVal.(int)
		if !ok || userID != 1 {
			helpers.ErrorJSON(w, errors.New("Require admin privilege"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ProvideUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			helpers.ErrorJSON(w, errors.New("User ID header missing"), http.StatusUnauthorized)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			helpers.ErrorJSON(w, errors.New("Invalid user ID format"), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
