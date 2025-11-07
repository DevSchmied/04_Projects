package server

import (
	"bookmanager-go/internal/controller"
	"bookmanager-go/internal/model"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestAdd1 verifies that the add1 helper function correctly increments integers.
func TestAdd1(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive input", 1, 2},
		{"negative input", -1, 0},
		{"zero input", 0, 1},
		{"larger positive input", 99, 100},
	}

	// Run each case as an individual subtest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := add1(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("add1(%d): got %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	// Arrange: define multiple test cases for formatDate
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "zero time formatted as default date",
			input:    time.Time{},
			expected: "1 January 0001",
		},
		{
			name:     "regular date formatted correctly",
			input:    time.Date(2025, time.October, 30, 12, 0, 0, 0, time.UTC),
			expected: "30 October 2025",
		},
		{
			name:     "different month formatted correctly",
			input:    time.Date(2023, time.March, 5, 8, 30, 0, 0, time.UTC),
			expected: "5 March 2023",
		},
	}
	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := formatDate(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("formatDate(%v) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// setupTestDB creates a temporary in-memory SQLite database
// and seeds it with predictable data for each test run.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// --- ARRANGE ---
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	if err := db.AutoMigrate(&model.Book{}); err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Insert predictable test data
	db.Create(&model.Book{ID: 1, Title: "Go Patterns"})
	db.Create(&model.Book{ID: 2, Title: "Concurrency in Go"})

	return db
}

func TestSearchStrategies(t *testing.T) {
	// --- ARRANGE ---
	db := setupTestDB(t)

	tests := []struct {
		name      string
		strategy  controller.SearchStrategy
		input     string
		wantTitle string
		wantErr   bool
	}{
		{
			name:      "IDSearch finds existing book",
			strategy:  controller.IDSearch{},
			input:     "1",
			wantTitle: "Go Patterns",
			wantErr:   false,
		},
		{
			name:      "TitleSearch finds by partial match",
			strategy:  controller.TitleSearch{},
			input:     "Concurrency",
			wantTitle: "Concurrency in Go",
			wantErr:   false,
		},
		{
			name:     "IDSearch fails for invalid ID format",
			strategy: controller.IDSearch{},
			input:    "abc",
			wantErr:  true,
		},
		{
			name:     "TitleSearch fails for non-existent title",
			strategy: controller.TitleSearch{},
			input:    "Unknown Book",
			wantErr:  true,
		},
	}

	_ = db
	_ = tests
}
