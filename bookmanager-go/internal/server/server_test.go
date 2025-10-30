package server

import (
	"testing"
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
