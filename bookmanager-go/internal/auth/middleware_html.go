package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRequiredHTML checks the JWT cookie for authentication and redirects to /login if invalid.
func AuthRequiredHTML() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read JWT from cookie
		tokenString, err := c.Cookie("jwt")
		if err != nil || tokenString == "" {
			// User not logged in → redirect to login
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

	}
}
