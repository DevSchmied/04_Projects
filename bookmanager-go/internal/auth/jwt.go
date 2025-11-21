package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: Replace hardcoded secret with an environment variable (e.g., from .env file)
var jwtSecret = []byte("SUPER_SECRET")

// CreateToken generates a JWT for a given user ID.
func CreateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken verifies and parses a JWT string.
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
