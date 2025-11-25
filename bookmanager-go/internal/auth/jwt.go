package auth

import (
	"bookmanager-go/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreateToken generates a JWT for a given user ID.
func CreateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}

// ValidateToken verifies and parses a JWT string.
func ValidateToken(tokenString string) (*jwt.Token, error) {
	secret := []byte(config.Cfg.JWTSecret)

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
