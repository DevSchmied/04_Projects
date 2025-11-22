package controller

import (
	"bookmanager-go/internal/model"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Initial BookController struct with DB dependency
type BookController struct {
	DB *gorm.DB
}

// findBookByParam chooses the right search strategy based on provided parameters.
func (bc *BookController) findBookByParam(idParam, title string) (*model.Book, error) {
	var (
		strategy SearchStrategy
		value    string
	)

	if idParam != "" {
		strategy = IDSearch{}
		value = idParam
	} else if title != "" {
		strategy = TitleSearch{}
		value = title
	} else {
		return nil, fmt.Errorf("no search parameter provided")
	}

	return strategy.Search(bc.DB, value)
}

// RegisterRoutes sets up all book-related routes.
func (bc *BookController) RegisterRoutes(r *gin.Engine) {

	r.GET("/", bc.ShowWelcomePage) // Display welcome page with logo
	r.GET("/list", bc.GetAllBooks) // Read all books

	adds := r.Group("/add")      // Group for add-related routes
	adds.GET("", bc.ShowAddPage) // Show the form to add a new book
	adds.POST("", bc.AddBook)    // Create a new book entry

	updates := r.Group("/update")                 // Group for update-related routes
	updates.GET("/search", bc.ShowSearchPage)     // Show the search form before updating
	updates.POST("/search", bc.FindBookForUpdate) // Search for a book to update
	updates.GET("/:id", bc.ShowEditPage)          // Show edit form for a specific book
	updates.POST("/:id", bc.UpdateBook)           // Update the selected book

	deletes := r.Group("/delete")                 // Group for delete-related routes
	deletes.GET("/search", bc.ShowSearchPage)     // Show the search form before deleting
	deletes.POST("/search", bc.FindBookForDelete) // Search for a book to delete
	deletes.POST("/:id", bc.DeleteBook)           // Delete the selected book

	r.GET("/:id", bc.GetBookByID) // Read a single book by its ID (must be last)
}
