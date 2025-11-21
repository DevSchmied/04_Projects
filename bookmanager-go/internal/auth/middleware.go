package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired returns a middleware that protects routes
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check header format
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		// Extract raw JWT string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token signature & expiration
		token, err := ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Parse claims from token (payload)
		claims := token.Claims.(jwt.MapClaims)

		// Convert "sub" claim (float64) to uint
		userID := uint(claims["sub"].(float64))

		// Save userID into Gin context (available for next handlers)
		c.Set("userID", userID)

		// Continue request processing
		c.Next()
	}
}
