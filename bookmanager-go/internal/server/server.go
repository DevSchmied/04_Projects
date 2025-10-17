package server

import (
	"bookmanager-go/internal/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server represents the HTTP server with its dependencies.
type Server struct {
	router         *gin.Engine
	bookController controller.BookController
	address        string
}

// NewServer creates a new Server instance with all dependencies injected.
func NewServer(db *gorm.DB, adr string) *Server {
	r := gin.Default()

	bc := controller.BookController{DB: db}

	return &Server{
		router:         r,
		bookController: bc,
		address:        adr,
	}
}

// StartWebServer starts the web server using dependency injection.
func (s *Server) Start() error {
	s.bookController.RegisterRoutes(s.router)

	// start the HTTP server
	return s.router.Run()
}
