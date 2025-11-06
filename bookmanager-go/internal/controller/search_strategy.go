package controller

import (
	"bookmanager-go/internal/model"
	"strconv"

	"gorm.io/gorm"
)

type SearchStrategy interface {
	Search(db *gorm.DB, value string) (*model.Book, error)
}

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
