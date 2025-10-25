package controller

import (
	"bookmanager-go/internal/model"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		books.GET("/", bc.GetAllBooks)                     // Read all
		books.GET("/:id", bc.GetBookByID)                  // Read one
		books.POST("/add", bc.AddBook)                     // Create
		books.POST("/update/search", bc.FindBookForUpdate) // Search before update
		books.POST("/update/:id", bc.UpdateBook)           // Update
		books.DELETE("/delete/:id", bc.DeleteBook)         // Delete
	}
}

// GetAllBooks renders an HTML page listing all books in the database.
func (bc *BookController) GetAllBooks(c *gin.Context) {
	var books []model.Book
	// Try to load all books from the database
	if err := bc.DB.Find(&books).Error; err != nil {
		log.Printf("Failed to retrieve books: %v\n", err)
		c.String(http.StatusInternalServerError, "Failed to retrieve books")
		return
	}
	// Render the list of books using an HTML template
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
		log.Printf("Invalid book ID format: %v\n", err)
		c.String(http.StatusBadRequest, "Invalid book ID format")
		return
	}
	// Ensure that the ID is greater than zero
	if id <= 0 {
		log.Printf("Invalid book ID value: %v\n", err)
		c.String(http.StatusBadRequest, "Invalid book ID value")
		return
	}

	// Retrieve book from database
	if err := bc.DB.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Book with ID %d not found\n", id)
			c.String(http.StatusNotFound, "Book not found")
		} else {
			log.Printf("Error fetching book with ID %d: %v\n", id, err)
			c.String(http.StatusInternalServerError, "Failed to fetch book")
		}
		return
	}

	// Render book details in an HTML template
	c.HTML(http.StatusOK, "book_details.html", gin.H{
		"Book": book,
	})
}

// AddBook handles HTML form submission to create a new book.
func (bc *BookController) AddBook(c *gin.Context) {
	// Read form input
	title := c.PostForm("title")
	author := c.PostForm("author")
	yearStr := c.PostForm("year")
	genre := c.PostForm("genre")
	isbn := c.PostForm("isbn")
	ratingStr := c.PostForm("rating")
	readStr := c.PostForm("read")

	// Validate required field
	if title == "" {
		c.HTML(http.StatusBadRequest, "book_add.html", gin.H{
			"error": "Title is required.",
		})
		return
	}

	// Convert numeric and boolean fields
	year, _ := strconv.Atoi(yearStr)
	rating, _ := strconv.ParseFloat(ratingStr, 64)

	// Round to one decimal place
	ratingStr = fmt.Sprintf("%.1f", rating)
	rating, _ = strconv.ParseFloat(ratingStr, 64)

	read := false
	readStr = strings.ToLower(readStr)
	if readStr == "yes" || readStr == "true" || readStr == "on" {
		read = true
	}

	// Create book object
	book := &model.Book{
		Title:  title,
		Author: author,
		Year:   year,
		Genre:  genre,
		ISBN:   isbn,
		Rating: rating,
		Read:   read,
	}

	// Save to database
	if err := bc.DB.Create(&book).Error; err != nil {
		log.Printf("Error saving book: %v\n", err)
		c.String(http.StatusInternalServerError, "Failed to save book")
		return
	}

	// Redirect to list page
	c.Redirect(http.StatusSeeOther, "/books")
}

// FindBookForUpdate searches a book by ID or title before updating it.
func (bc *BookController) FindBookForUpdate(c *gin.Context) {
	idParam := c.Param("id")
	title := c.Param("title")

	// Validate that at least one parameter is provided
	if idParam == "" && title == "" {
		log.Println("ID or title is required.")
		c.HTML(http.StatusBadRequest, "book_search.html", gin.H{
			"error": "Please provide either ID or title to search.",
		})
		return
	}

	var (
		id  int
		err error
	)

	// Convert ID from string to integer (only if provided)
	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			log.Printf("Invalid book ID format: %v\n", err)
			c.String(http.StatusBadRequest, "Invalid book ID format (must be a number)")
			return
		}
	}

	var book model.Book
	// Search for a book by ID or title in the database
	if err := bc.DB.Where("id = ? OR title = ?", id, title).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Book not found: %v\n", err)
			c.String(http.StatusNotFound, "Book not found")
			return
		}
		log.Printf("Error while searching for book: %v\n", err)
		c.String(http.StatusInternalServerError, "Internal error while searching for book")
		return
	}

	// Render the edit page with the found book data
	c.HTML(http.StatusOK, "book_edit.html", gin.H{
		"Book": book,
	})
}

// UpdateBook updates an existing book in the database based on form input.
func (bc *BookController) UpdateBook(c *gin.Context) {
	// Read form input
	idParam := c.PostForm("id")
	title := c.PostForm("title")
	author := c.PostForm("author")
	yearStr := c.PostForm("year")
	genre := c.PostForm("genre")
	isbn := c.PostForm("isbn")
	ratingStr := c.PostForm("rating")
	readStr := c.PostForm("read")

	// Validate required field
	if title == "" {
		c.HTML(http.StatusBadRequest, "book_edit.html", gin.H{
			"error": "Title is required.",
		})
		return
	}

	// Convert numeric and boolean fields
	id, _ := strconv.Atoi(idParam)
	year, _ := strconv.Atoi(yearStr)
	rating, _ := strconv.ParseFloat(ratingStr, 64)

	// Round to one decimal place
	ratingStr = fmt.Sprintf("%.1f", rating)
	rating, _ = strconv.ParseFloat(ratingStr, 64)

	// Normalize read flag
	read := false
	readStr = strings.ToLower(readStr)
	if readStr == "yes" || readStr == "true" || readStr == "on" {
		read = true
	}

	// Create book object
	book := &model.Book{
		ID:     uint(id),
		Title:  title,
		Author: author,
		Year:   year,
		Genre:  genre,
		ISBN:   isbn,
		Rating: rating,
		Read:   read,
	}

	// Update database record
	if err := bc.DB.Model(&model.Book{}).Where("id = ?", id).Updates(book).Error; err != nil {
		log.Printf("Error updating book ID %d: %v\n", id, err)
		c.String(http.StatusInternalServerError, "Failed to update book")
		return
	}

	// Redirect on success
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/books/%d", id))
}

// FindBookForDelete searches for a book by ID or title before deletion.
func (bc *BookController) FindBookForDelete(c *gin.Context) {
	idParam := c.Param("id")
	_ = idParam
	title := c.Param("title")
	_ = title

	// Validate that at least one parameter is provided
	if idParam == "" && title == "" {
		log.Println("ID or title is required.")
		c.HTML(http.StatusBadRequest, "book_search.html", gin.H{
			"error": "Please provide either ID or title to search.",
		})
		return
	}

	var (
		id  int
		err error
	)
	_ = id

	// Convert ID from string to integer (only if provided)
	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			log.Printf("Invalid book ID format: %v\n", err)
			c.String(http.StatusBadRequest, "Invalid book ID format (must be a number)")
			return
		}
	}

	var book model.Book
	// Search for a book by ID or title in the database
	if err := bc.DB.Where("id = ? OR title = ?", id, title).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Book not found: %v\n", err)
			c.String(http.StatusNotFound, "Book not found")
			return
		}
		log.Printf("Error while searching for book: %v\n", err)
		c.String(http.StatusInternalServerError, "Internal error while searching for book")
		return
	}

	// Render the delete confirmation page
	c.HTML(http.StatusOK, "book_delete.html", gin.H{
		"Book": book,
	})
}

func (bc *BookController) DeleteBook(c *gin.Context) {}
