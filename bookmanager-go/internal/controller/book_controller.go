package controller

import (
	"bookmanager-go/internal/model"
	"errors"
	"net/http"
	"strconv"

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

// GetBookByID returns a single book by its ID as JSON.
func (bc *BookController) GetBookByID(c *gin.Context) {
	idParam := c.Param("id")
	var book model.Book

	// Prove that the ID is a valid integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	// Ensure that the ID is greater than zero
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := bc.DB.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, book)
}

func (bc *BookController) AddBook(c *gin.Context)    {}
func (bc *BookController) UpdateBook(c *gin.Context) {}
func (bc *BookController) DeleteBook(c *gin.Context) {}
