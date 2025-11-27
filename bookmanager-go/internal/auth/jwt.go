package auth

import (
	"bookmanager-go/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService kapselt alle JWT-bezogenen Operationen (Signieren, Validieren)
type JWTService struct {
	Secret string
}

// NewJWTService erstellt einen neuen JWTService auf Basis der geladenen AppConfig.
func NewJWTService(cfg *config.AppConfig) *JWTService {
	return &JWTService{Secret: cfg.JWTSecret}
}

// CreateToken generates a JWT for a given user ID.
func (j *JWTService) CreateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}

// ValidateToken verifies and parses a JWT string.
func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
}
