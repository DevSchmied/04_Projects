package controller

import (
	"bookmanager-go/internal/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ShowAddPage displays the HTML form to add a new book.
func (bc *BookController) ShowAddPage(c *gin.Context) {
	bc.renderHTML(c, http.StatusOK, "book_add.html", gin.H{
		"Title":       "Add Book",
		"PageTitle":   "Add a New Book",
		"Description": "Enter the details of the new book and click Save.",
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
		bc.renderHTML(c, http.StatusBadRequest, "book_add.html", gin.H{
			"Title":       "Add Book",
			"Message":     "Title is required.",
			"MessageType": "warning",
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
		bc.renderHTML(c, http.StatusInternalServerError, "book_add.html", gin.H{
			"Title":       "Add Book",
			"Message":     "Failed to save book.",
			"MessageType": "danger",
		})
		return
	}

	if bc.Cache != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Millisecond)
		defer cancel()

		if err := bc.Cache.InvalidateBookList(ctx); err != nil {
			log.Printf("Cache invalidate error: %v", err)
		} else {
			log.Println("Cache invalidated: book list cleared")
		}
	}

	// Redirect to list page
	c.Redirect(http.StatusSeeOther, "/books/list")
}
