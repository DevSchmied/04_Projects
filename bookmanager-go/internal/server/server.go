package server

import (
	"bookmanager-go/internal/controller"
	"text/template"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server represents the HTTP server with its dependencies.
type Server struct {
	router         *gin.Engine
	bookController controller.BookController
	address        string
	staticRoute    string
	staticPath     string
	templatePath   string
}

// NewServer creates a new Server instance with all dependencies injected.
func NewServer(db *gorm.DB, adr, templates, staticRoute, staticPath string) *Server {
	r := gin.Default()

	bc := controller.BookController{DB: db}

	return &Server{
		router:         r,
		bookController: bc,
		address:        adr,
		staticRoute:    staticRoute,
		staticPath:     staticPath,
		templatePath:   templates,
	}
}

// StartWebServer starts the web server using dependency injection.
func (s *Server) Start() error {

	// Serve static files
	s.router.Static(s.staticRoute, s.staticPath)

	// Register custom template functions
	funcMap := template.FuncMap{
		"add1": add1,
	}

	// Parse templates with custom functions
	tmpl := template.Must(
		template.New("base").
			Funcs(funcMap).
			ParseGlob(s.templatePath),
	)

	_ = tmpl

	// Load all HTML templates
	s.router.LoadHTMLGlob(s.templatePath)

	// Register routes for controllers
	s.bookController.RegisterRoutes(s.router)

	// start the HTTP server
	return s.router.Run(s.address)
}

// add1 adds 1 to an integer index (used in templates for human-friendly numbering)
func add1(i int) int {
	return i + 1
}
