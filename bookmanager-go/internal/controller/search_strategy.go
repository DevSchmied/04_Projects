package controller

import (
	"bookmanager-go/internal/model"

	"gorm.io/gorm"
)

type SearchStrategy interface {
	Search(db *gorm.DB, value string) (*model.Book, error)
}
