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

// renderHTML is a helper function that renders an HTML template
func (bc *BookController) renderHTML(c *gin.Context, status int, templateName string, data gin.H) {
	// Ensure the map is initialized
	if data == nil {
		data = gin.H{}
	}

	// Global default values
	if _, exists := data["Title"]; !exists {
		data["Title"] = "BookManager"
	}
	if _, exists := data["PageTitle"]; !exists {
		if status >= 400 {
			data["PageTitle"] = "Error Situation"
		} else {
			data["PageTitle"] = "BookManager"
		}
	}
	if _, exists := data["Description"]; !exists {
		if status >= 400 {
			data["Description"] = "An unexpected error occurred. Please try again later."
		} else {
			data["Description"] = "Manage your personal library — add, edit, and organize your favorite books."
		}
	}
	if _, exists := data["MessageType"]; !exists {
		data["MessageType"] = ""
	}
	if _, exists := data["Message"]; !exists {
		data["Message"] = ""
	}

	// Render page
	c.HTML(status, templateName, data)
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

// parseIDParam extracts and validates the "id" URL parameter.
func parseIDParam(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	if idParam == "" {
		return 0, fmt.Errorf("missing ID parameter")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid ID: %s", idParam)
	}
	return uint(id), nil
}

// RegisterRoutes sets up all book-related routes.
func (bc *BookController) RegisterRoutes(r *gin.Engine) {
	books := r.Group("/books")
	{
		books.GET("/", bc.ShowWelcomePage) // Display welcome page with logo
		books.GET("/list", bc.GetAllBooks) // Read all books
		books.GET("/add", bc.ShowAddPage)  // Show the form to add a new book
		books.POST("/add", bc.AddBook)     // Create a new book entry

		books.GET("/update/search", bc.ShowSearchPage)     // Show the search form before updating
		books.POST("/update/search", bc.FindBookForUpdate) // Search for a book to update
		books.GET("/update/:id", bc.ShowEditPage)          // Show edit form for a specific book
		books.POST("/update/:id", bc.UpdateBook)           // Update the selected book

		books.GET("/delete/search", bc.ShowSearchPage)     // Show the search form before deleting
		books.POST("/delete/search", bc.FindBookForDelete) // Search for a book to delete
		books.POST("/delete/:id", bc.DeleteBook)           // Delete the selected book

		books.GET("/:id", bc.GetBookByID) // Read a single book by its ID (must be last)
	}
}

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

	// Redirect to list page
	c.Redirect(http.StatusSeeOther, "/books/list")
}

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
			})
		} else {
			bc.renderHTML(c, http.StatusInternalServerError, "book_search.html", gin.H{
				"Title":       "Book Search",
				"Message":     "Failed to search for book.",
				"MessageType": "danger",
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

	// Success message after update
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
	idParam := c.Param("id")

	var (
		id  int
		err error
	)

	// Validate and convert ID
	if idParam == "" {
		log.Println("ID is required for deletion")
		bc.renderHTML(c, http.StatusBadRequest, "books_list.html", gin.H{
			"Title":       "Book List",
			"Message":     "Book ID is required for deletion.",
			"MessageType": "warning",
		})
		return
	}

	id, err = strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid ID format: %v\n", err)
		bc.renderHTML(c, http.StatusBadRequest, "books_list.html", gin.H{
			"Title":       "Book List",
			"Message":     fmt.Sprintf("Invalid ID format: %v", err),
			"MessageType": "danger",
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
