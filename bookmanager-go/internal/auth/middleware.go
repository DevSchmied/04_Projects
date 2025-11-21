package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired returns a middleware that protects routes
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if header exists and starts with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		// Extract the token string (remove "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT token
		token, err := ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Pass control to next handler
		c.Next()
	}
}
