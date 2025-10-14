package controller

import (
	"bookmanager-go/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Initial BookController struct with DB dependency
type BookController struct {
	DB *gorm.DB
}

// RegisterRoutes sets up all book-related routes.
func (bc *BookController) RegisterRoutes(r *gin.Engine) {
	books := r.Group("/books")
	{
		books.GET("/", bc.GetAllBooks)             // Read all
		books.GET("/:id", bc.GetBookByID)          // Read one
		books.POST("/add", bc.AddBook)             // Create
		books.PUT("/:id/edit", bc.UpdateBook)      // Update
		books.DELETE("/delete/:id", bc.DeleteBook) // Delete
	}
}

// GetAllBooks returns all books as JSON.
func (bc *BookController) GetAllBooks(c *gin.Context) {
	var books []model.Book
	if result := bc.DB.Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (bc *BookController) GetBookByID(c *gin.Context) {}
func (bc *BookController) AddBook(c *gin.Context)     {}
func (bc *BookController) UpdateBook(c *gin.Context)  {}
func (bc *BookController) DeleteBook(c *gin.Context)  {}
