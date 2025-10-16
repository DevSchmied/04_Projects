package server

import (
	"bookmanager-go/internal/controller"

	"github.com/gin-gonic/gin"
)

// StartWebServer2 starts the web server using dependency injection.
func StartWebServer(r *gin.Engine, bc controller.BookController, address string) error {
	bc.RegisterRoutes(r)

	// start the HTTP server
	return r.Run(address)
}
