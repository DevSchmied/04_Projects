package server

import (
	"bookmanager-go/internal/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StartWebServer initializes the Gin router, registers routes, and starts the server.
func StartWebServer(db *gorm.DB) {
	router := gin.Default()

	bookController := controller.BookController{DB: db}
	bookController.RegisterRoutes(router) // register all /books routes

	// start the HTTP server on localhost:8080
	router.Run("localhost:8080")
}
