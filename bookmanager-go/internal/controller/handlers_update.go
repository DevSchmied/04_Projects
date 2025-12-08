package controller

import (
	"bookmanager-go/internal/model"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShowEditPage loads the edit form pre-filled with book data
func (bc *BookController) ShowEditPage(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		log.Printf("Invalid book ID: %v\n", err)
		bc.renderHTML(c, http.StatusBadRequest, "books_list.html", gin.H{
			"Title":       "Edit Book",
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
				"Title":       "Edit Book",
				"PageTitle":   "Edit Book",
				"Description": "Update the book information and save your changes.",
				"Message":     "Book not found",
				"MessageType": "info",
				"Book":        book,
			})
			return
		} else {
			log.Printf("Error fetching book with ID %d: %v\n", id, err)
			bc.renderHTML(c, http.StatusInternalServerError, "books_list.html", gin.H{
				"Title":       "Edit Book",
				"Message":     "An unexpected database error occurred.",
				"MessageType": "danger",
				"Book":        book,
			})
		}
		return
	}

	// Render book details in an HTML template
	bc.renderHTML(c, http.StatusOK, "book_edit.html", gin.H{
		"Title":       "Edit Book",
		"PageTitle":   fmt.Sprintf("Edit: %s", book.Title),
		"Description": "Update the book information and save your changes.",
		"Book":        book,
	})
}

// FindBookForUpdate searches a book by ID or title before updating it.
func (bc *BookController) FindBookForUpdate(c *gin.Context) {
	idParam := c.PostForm("id")
	title := c.PostForm("title")

	// Validate that at least one parameter is provided
	if idParam == "" && title == "" {
		log.Println("ID or title is required.")
		bc.renderHTML(c, http.StatusOK, "book_search.html", gin.H{
			"Title":       "Book Search",
			"PageTitle":   "Book Search",
			"Description": "Enter either the book ID or title to search for a specific book.",
			"Message":     "Please provide either ID or title to search.",
			"MessageType": "info",
			"Action":      "update/search",
		})
		return
	}

	// Use helper function to find the book
	book, err := bc.findBookByParam(idParam, title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			bc.renderHTML(c, http.StatusNotFound, "book_search.html", gin.H{
				"Title":       "Book Search",
				"Message":     "Book not found.",
				"MessageType": "info",
				"Action":      "update/search",
			})
		} else {
			bc.renderHTML(c, http.StatusInternalServerError, "book_search.html", gin.H{
				"Title":       "Book Search",
				"Message":     "Failed to search for book.",
				"MessageType": "danger",
				"Action":      "update/search",
			})
		}
		return
	}

	// Render the edit page with the found book data
	bc.renderHTML(c, http.StatusOK, "book_edit.html", gin.H{
		"Title":       "Book Edit",
		"PageTitle":   fmt.Sprintf("Edit: %s", book.Title),
		"Description": "Update the book information and save your changes.",
		"Book":        book,
	})
}

// UpdateBook updates an existing book in the database based on form input.
func (bc *BookController) UpdateBook(c *gin.Context) {
	// Read form input
	idParam := c.Param("id")
	title := c.PostForm("title")
	author := c.PostForm("author")
	yearStr := c.PostForm("year")
	genre := c.PostForm("genre")
	isbn := c.PostForm("isbn")
	ratingStr := c.PostForm("rating")
	readStr := c.PostForm("read")

	// Validate required field
	if title == "" {
		bc.renderHTML(c, http.StatusBadRequest, "book_edit.html", gin.H{
			"Title":       "Edit Book",
			"Message":     "Title is required.",
			"MessageType": "warning",
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
		bc.renderHTML(c, http.StatusInternalServerError, "book_edit.html", gin.H{
			"Title":       "Edit Book",
			"Message":     fmt.Sprintf("Failed to update book ID %d.", id),
			"MessageType": "danger",
			"Book":        book,
		})
		return
	}

	if bc.cache != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Millisecond)
		defer cancel()

		if err := bc.cache.InvalidateBookList(ctx); err != nil {
			log.Printf("Cache invalidate error: %v", err)
		}
	}

	log.Printf("Book with ID %d successfully updated.\n", id)
	bc.renderHTML(c, http.StatusOK, "book_edit.html", gin.H{
		"Title":       "Edit Book",
		"PageTitle":   "Edit Book",
		"Description": "Update the book information and save your changes.",
		"Message":     fmt.Sprintf("Book with ID %d was successfully updated.", id),
		"MessageType": "success",
		"Book":        book,
	})
}
