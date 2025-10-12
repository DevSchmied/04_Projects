package service

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// --- Mock Implementations ---

// mockConnector simulates different database connection outcomes.
type mockConnector struct {
	returnDB  *gorm.DB
	returnErr error
}

// Connect returns either a mock *gorm.DB or an error, depending on the test setup.
func (m *mockConnector) Connect() (*gorm.DB, error) {
	return m.returnDB, m.returnErr
}

// --- Table-driven Tests ---
// TestInitDB follows the Arrange–Act–Assert (AAA) pattern for clarity.
func TestInitDB(t *testing.T) {
	// --- ARRANGE ---
	// Define the table of test cases with different expected behaviors.
	tests := []struct {
		name      string
		connector DatabaseConnector
		wantPanic bool
	}{
		{
			name:      "successful connection",
			connector: &mockConnector{returnDB: &gorm.DB{}, returnErr: nil},
			wantPanic: false,
		},
		{
			name:      "failed connection",
			connector: &mockConnector{returnDB: nil, returnErr: errors.New("failed")},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// --- ARRANGE ---
			// Prepare the connector and expected panic recovery if needed.
			connector := tt.connector

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic but none occurred")
					}
				}()
			}

			// --- ACT ---
			// Execute the function under test.
			db := InitDB(connector)

			// --- ASSERT ---
			// For successful connections, ensure db is not nil.
			if !tt.wantPanic && db == nil {
				t.Errorf("expected *gorm.DB instance, got nil")
			}
		})
	}
}
