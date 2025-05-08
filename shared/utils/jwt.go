package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID int, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
