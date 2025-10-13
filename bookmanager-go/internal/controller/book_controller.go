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

// GetAllBooks returns all books as JSON.
func (bc *BookController) GetAllBooks(c *gin.Context) {
	var books []model.Book
	if result := bc.DB.Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}
