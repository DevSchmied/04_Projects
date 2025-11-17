package controller

import (
	"bookmanager-go/internal/model"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShowWelcomePage renders the welcome page.
func (bc *BookController) ShowWelcomePage(c *gin.Context) {
	bc.renderHTML(c, http.StatusOK, "index.html", gin.H{
		"Title":       "Welcome to BookManager",
		"PageTitle":   "Welcome to BookManager",
		"Description": "Manage your personal library — add, edit, and organize your favorite books.",
	})
}

// GetAllBooks renders an HTML page listing all books in the database.
func (bc *BookController) GetAllBooks(c *gin.Context) {
	var books []model.Book
	// Try to load all books from the database
	if err := bc.DB.Find(&books).Error; err != nil {
		log.Printf("Failed to retrieve books: %v\n", err)
		bc.renderHTML(c, http.StatusInternalServerError, "books_list.html", gin.H{
			"Title":       "Book List",
			"Message":     "Failed to retrieve books",
			"MessageType": "warning",
			"Books":       books,
		})
		return
	}

	// Show info message if no books are found
	if len(books) == 0 {
		bc.renderHTML(c, http.StatusOK, "books_list.html", gin.H{
			"Title":       "Book List",
			"PageTitle":   "Your Library",
			"Description": "List of all books currently stored in your library.",
			"Message":     "There are no saved books yet. Add your first book!",
			"MessageType": "info",
			"Books":       books,
		})
		return
	}

	// Render the list of books
	bc.renderHTML(c, http.StatusOK, "books_list.html", gin.H{
		"Title":       "Book List",
		"Description": "List of all books currently stored in your library.",
		"Books":       books,
	})
}

// GetBookByID renders a page showing details of a single book by its ID.
func (bc *BookController) GetBookByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		log.Printf("Invalid book ID: %v\n", err)
		bc.renderHTML(c, http.StatusBadRequest, "books_list.html", gin.H{
			"Title":       "Book Details",
			"Message":     err.Error(),
			"MessageType": "warning",
		})
		return
	}

	var book model.Book

	// Retrieve book from database
	if err := bc.DB.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Book with ID %d not found\n", id)
			bc.renderHTML(c, http.StatusNotFound, "books_list.html", gin.H{
				"Title":       "Book Details",
				"Message":     "Book not found.",
				"MessageType": "info",
				"Book":        book,
			})
		} else {
			log.Printf("Error fetching book with ID %d: %v\n", id, err)
			bc.renderHTML(c, http.StatusInternalServerError, "books_list.html", gin.H{
				"Title":       "Book Details",
				"Message":     "An unexpected database error occurred.",
				"MessageType": "danger",
				"Book":        book,
			})
		}
		return
	}

	// Render book details in an HTML template
	bc.renderHTML(c, http.StatusOK, "book_details.html", gin.H{
		"Title":       "Book Details",
		"PageTitle":   fmt.Sprintf("Details of '%s'", book.Title),
		"Description": "Detailed information about the selected book from your library.",
		"Book":        book,
	})
}

// ShowSearchPage renders HTML search form
func (bc *BookController) ShowSearchPage(c *gin.Context) {
	path := c.FullPath()

	var pageTitle, description, action string

	if strings.Contains(path, "update") {
		pageTitle = "Find Book to Update"
		description = "Enter either the book ID or title to search for a book you want to update."
		action = "update/search"
	} else if strings.Contains(path, "delete") {
		pageTitle = "Find Book to Delete"
		description = "Enter either the book ID or title to search for a book you want to delete."
		action = "delete/search"
	} else {
		pageTitle = "Book Search"
		description = "Enter either the book ID or title to search for a specific book."
		action = "update/search"
	}

	bc.renderHTML(c, http.StatusOK, "book_search.html", gin.H{
		"Title":       pageTitle,
		"PageTitle":   pageTitle,
		"Description": description,
		"Action":      action,
	})
}
