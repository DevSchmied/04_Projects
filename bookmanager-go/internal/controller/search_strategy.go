package controller

import (
	"bookmanager-go/internal/model"
	"strconv"

	"gorm.io/gorm"
)

// SearchStrategy defines a common interface for different search strategies.
type SearchStrategy interface {
	Search(db *gorm.DB, value string) (*model.Book, error)
}

// IDSearch implements SearchStrategy for finding a book by its numeric ID.
type IDSearch struct{}

func (s IDSearch) Search(db *gorm.DB, value string) (*model.Book, error) {
	var book model.Book
	id, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	err = db.First(&book, id).Error
	return &book, err
}

// TitleSearch implements SearchStrategy for finding a book by its title.
type TitleSearch struct{}

func (s TitleSearch) Search(db *gorm.DB, value string) (*model.Book, error) {
	var book model.Book
	err := db.Where("title LIKE ?", "%"+value+"%").First(&book).Error
	return &book, err
}
