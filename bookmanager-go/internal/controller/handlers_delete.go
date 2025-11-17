package controller

import (
	"bookmanager-go/internal/model"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FindBookForDelete searches for a book by ID or title before deletion.
func (bc *BookController) FindBookForDelete(c *gin.Context) {
	idParam := c.PostForm("id")
	title := c.PostForm("title")

	// Validate that at least one parameter is provided
	if idParam == "" && title == "" {
		log.Println("ID or title is required.")
		bc.renderHTML(c, http.StatusBadRequest, "book_search.html", gin.H{
			"Title":       "Find Book to Delete",
			"Message":     "Please provide either the book ID or title.",
			"MessageType": "warning",
			"Action":      "delete/search",
		})
		return
	}

	// Use helper function to find the book
	book, err := bc.findBookByParam(idParam, title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Book not found for ID '%s' or title '%s'\n", idParam, title)
			bc.renderHTML(c, http.StatusNotFound, "book_search.html", gin.H{
				"Title":       "Find Book to Delete",
				"Message":     "Book not found.",
				"MessageType": "info",
				"Action":      "delete/search",
			})
			return
		}

		log.Printf("Error searching for book: %v\n", err)
		bc.renderHTML(c, http.StatusInternalServerError, "book_search.html", gin.H{
			"Title":       "Find Book to Delete",
			"Message":     "An unexpected error occurred while searching for the book.",
			"MessageType": "danger",
			"Action":      "delete/search",
		})
		return
	}

	// Render the delete confirmation page
	bc.renderHTML(c, http.StatusOK, "book_delete.html", gin.H{
		"Title":       "Confirm Book Deletion",
		"PageTitle":   "Confirm Book Deletion",
		"Description": "Please confirm that you want to permanently delete this book from your library.",
		"Book":        book,
	})
}

// DeleteBook deletes a book from the database based on its ID.
func (bc *BookController) DeleteBook(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		log.Printf("Invalid book ID: %v\n", err)
		bc.renderHTML(c, http.StatusBadRequest, "books_list.html", gin.H{
			"Title":       "Book List",
			"Message":     err.Error(),
			"MessageType": "warning",
		})
		return
	}

	// Check if the book exists
	var book model.Book
	if err = bc.DB.First(&book, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Book with ID %d not found\n", id)
			bc.renderHTML(c, http.StatusNotFound, "books_list.html", gin.H{
				"Title":       "Book List",
				"Message":     fmt.Sprintf("Book with ID %d not found.", id),
				"MessageType": "info",
			})
			return
		}

		log.Printf("Error fetching book before deletion: %v\n", err)
		bc.renderHTML(c, http.StatusInternalServerError, "books_list.html", gin.H{
			"Title":       "Book List",
			"Message":     "An unexpected error occurred while checking the book.",
			"MessageType": "danger",
		})
		return
	}

	// Delete book record from database
	if err = bc.DB.Delete(&book).Error; err != nil {
		log.Printf("Error deleting book ID %d: %v\n", id, err)
		bc.renderHTML(c, http.StatusInternalServerError, "books_list.html", gin.H{
			"Title":       "Book List",
			"Message":     "Failed to delete the book.",
			"MessageType": "danger",
		})
		return
	}

	// Redirect back to book list with success info
	log.Printf("Book with ID %d successfully deleted.\n", id)
	bc.renderHTML(c, http.StatusOK, "books_list.html", gin.H{
		"Title":       "Book List",
		"PageTitle":   "Your Library",
		"Description": "List of all books currently stored in your library.",
		"Message":     fmt.Sprintf("Book with ID %d was successfully deleted.", id),
		"MessageType": "success",
	})
}
