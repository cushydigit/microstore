package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cushydigit/microstore/shared/helpers"
	"github.com/cushydigit/microstore/shared/types"
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

		claims := &types.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			helpers.ErrorJSON(w, errors.New("Invalid or expired token"), http.StatusUnauthorized)
			return
		}

		// inject into header for downstream services
		r.Header.Set(string(types.XUserID), strconv.Itoa(claims.UserID))
		r.Header.Set(string(types.XUserEmail), claims.Email)

		// inject into context for internal use
		ctx := context.WithValue(r.Context(), types.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, types.UserEmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emailVal := r.Context().Value(types.UserEmailKey)
		email, ok := emailVal.(string)
		if !ok || !strings.HasSuffix(email, "@admin.microstore.com") {
			helpers.ErrorJSON(w, errors.New("Require admin privilege"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ProvideUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get(string(types.XUserID))
		if userIDStr == "" {
			helpers.ErrorJSON(w, errors.New("User ID header missing"), http.StatusUnauthorized)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			helpers.ErrorJSON(w, errors.New("Invalid user ID format"), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), types.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
