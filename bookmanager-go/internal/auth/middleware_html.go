package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequiredHTML checks the JWT cookie for authentication and redirects to /login if invalid.
func AuthRequiredHTML(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read JWT from cookie
		tokenString, err := c.Cookie("jwt")
		if err != nil || tokenString == "" {
			// User not logged in → redirect to login
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Set("userID", uint(claims["sub"].(float64)))
		c.Next()
	}
}
