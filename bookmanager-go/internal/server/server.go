package server

import (
	"bookmanager-go/internal/auth"
	"bookmanager-go/internal/controller"
	"html/template"
	"time"

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

// setupTemplates registers custom template functions and loads all templates.
func (s *Server) setupTemplates() {
	// Register custom template functions
	funcsMap := template.FuncMap{
		"add1":       add1,
		"formatDate": formatDate,
	}

	// Parse templates with custom functions
	tmpl := template.Must(
		template.New("base").
			Funcs(funcsMap).
			ParseGlob(s.templatePath),
	)

	// Set the parsed template for Gin
	s.router.SetHTMLTemplate(tmpl)
}

// Start starts the web server using dependency injection.
func (s *Server) Start() error {
	// Serve static files
	s.router.Static(s.staticRoute, s.staticPath)

	// Register custom template functions and load all templates.
	s.setupTemplates()

	s.registerPublicRoutes()    // Register routes that are accessible without authentication
	s.registerProtectedRoutes() // Register routes that require a valid JWT cookie

	// start the HTTP server
	return s.router.Run(s.address)
}

// add1 adds 1 to an integer index (used in templates for human-friendly numbering)
func add1(i int) int {
	return i + 1
}

// formatDate formats a time.Time value into a human-readable string.
func formatDate(d time.Time) string {
	return d.Format("2 January 2006")
}

// registerPublicRoutes defines all routes that do not require authentication.
func (s *Server) registerPublicRoutes() {
	// Initialize HTML authentication controller
	authHTML := controller.AuthHTMLController{DB: s.bookController.DB}

	// Public registration routes
	s.router.GET("/register", authHTML.ShowRegisterPage)
	s.router.POST("/register", authHTML.RegisterUser)

	// Public login routes
	s.router.GET("/login", authHTML.ShowLoginPage)
	s.router.POST("/login", authHTML.LoginUser)

	// Public logout route
	s.router.GET("/logout", authHTML.LogoutUser)
}

// registerProtectedRoutes defines all routes that require authentication.
func (s *Server) registerProtectedRoutes() {
	// Group all /books routes under a protected route group
	books := s.router.Group("/books")

	// Apply HTML JWT authentication middleware
	books.Use(auth.AuthRequiredHTML())

	// Register routes for controllers
	s.bookController.RegisterRoutes(books)
}
