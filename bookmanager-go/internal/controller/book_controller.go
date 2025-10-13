package controller

import (
	"gorm.io/gorm"
)

// initial BookController struct with DB dependency
type BookController struct {
	DB *gorm.DB
}
