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

// GetAllBooks renders an HTML page listing all books in the database.
func (bc *BookController) GetAllBooks(c *gin.Context) {
	var books []model.Book
	if err := bc.DB.Find(&books).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve books")
		return
	}
	c.HTML(http.StatusOK, "books_list.html", gin.H{
		"Books": books,
	})
}

// GetBookByID renders a page showing details of a single book by its ID.
func (bc *BookController) GetBookByID(c *gin.Context) {
	idParam := c.Param("id")
	var book model.Book

	// Prove that the ID is a valid integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid book ID format")
		return
	}
	// Ensure that the ID is greater than zero
	if id <= 0 {
		c.String(http.StatusBadRequest, "Invalid book ID value")
		return
	}

	if err := bc.DB.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, "Book not found")
		} else {
			c.String(http.StatusInternalServerError, "Failed to fetch book")
		}
		return
	}
	c.HTML(http.StatusOK, "book_details.html", gin.H{
		"Book": book,
	})
}

func (bc *BookController) AddBook(c *gin.Context)    {}
func (bc *BookController) UpdateBook(c *gin.Context) {}
func (bc *BookController) DeleteBook(c *gin.Context) {}
