package model

import (
	"time"

	"gorm.io/gorm"
)

// Book represents a book entity stored in the database.

// Note:
// This is a denormalized (non-normalized) model structure.
// It is sufficient for this phase of the project (CRUD and testing).
// In later stages, the model can be normalized.
type Book struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Author    string         `json:"author"`
	Year      int            `json:"year"`
	Genre     string         `json:"genre"`
	ISBN      string         `json:"isbn"`
	Rating    float64        `json:"rating"`
	Read      bool           `json:"read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}
