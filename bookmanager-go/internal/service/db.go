package service

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Interface for database connections.
type DatabaseConnector interface {
	Connect() (*gorm.DB, error)
}

// Concrete implementation of DatabaseConnector for SQLite.
type SQLiteConnector struct {
	DBPath string
}

// Establishes a connection to the SQLite database file.
func (s *SQLiteConnector) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(s.DBPath), &gorm.Config{})
	return db, err
}

// High-level initialization function (supports dependency injection).
func InitDB(d DatabaseConnector) *gorm.DB {
	db, err := d.Connect()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	log.Println("Success: connected to database.")
	return db
}
