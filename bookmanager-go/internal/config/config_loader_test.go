package config

import (
	"errors"
	"fmt"
	"testing"
)

// FakeReader implements EnvReader for testing.
type FakeReader struct {
	Values map[string]string
}

func (f FakeReader) Get(key string) string {
	return f.Values[key]
}

// Test: Successful load

func TestLoadConfig_Success(t *testing.T) {
	// Arrange
	fakeReader := FakeReader{
		Values: map[string]string{
			"JWT_SECRET": "test-secret",
		},
	}

	loadCount := 0
	fakeLoad := func() error {
		loadCount++
		return nil
	}

	loader := NewEnvLoader("ignored", fakeLoad, fakeReader)

	// Act
	cfg, err := loader.LoadConfig(0, 3)

	// Assert
	if err != nil {
		t.Fatalf("LoadConfig returned unexpected error: %v", err)
	}
	if cfg.JWTSecret != "test-secret" {
		t.Errorf("expected JWT_SECRET = 'test-secret', got '%s'", cfg.JWTSecret)
	}
	if loadCount != 1 {
		t.Errorf("expected Load() to be called once, got %d calls", loadCount)
	}
}

// Test: Retry logic works (first fail, then success)
func TestLoadConfig_RetrySuccess(t *testing.T) {
	// Arrange
	fakeReader := FakeReader{
		Values: map[string]string{
			"JWT_SECRET": "retry-secret",
		},
	}

	loadCount := 0
	fakeLoad := func() error {
		loadCount++
		if loadCount == 1 {
			return errors.New("first attempt fails")
		}
		return nil
	}

	loader := NewEnvLoader("ignored", fakeLoad, fakeReader)

	// Act
	cfg, err := loader.LoadConfig(0, 3)

	// Assert
	if err != nil {
		t.Fatalf("LoadConfig returned unexpected error: %v", err)
	}
	if cfg.JWTSecret != "retry-secret" {
		t.Errorf("expected JWT_SECRET 'retry-secret', got '%s'", cfg.JWTSecret)
	}
	if loadCount != 2 {
		t.Errorf("expected 2 attempts, got %d", loadCount)
	}
}

// Test: Max attempts reached (always failing)

func TestLoadConfig_MaxAttemptsFailure(t *testing.T) {
	// Arrange
	fakeReader := FakeReader{Values: map[string]string{}}

	loadCount := 0
	fakeLoad := func() error {
		loadCount++
		return fmt.Errorf("fail %d", loadCount)
	}

	loader := NewEnvLoader("ignored", fakeLoad, fakeReader)

	// Act
	cfg, err := loader.LoadConfig(0, 3)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if cfg != nil {
		t.Errorf("expected cfg to be nil on failure, got %+v", cfg)
	}
	if loadCount != 3 {
		t.Errorf("expected load to be called 3 times, got %d", loadCount)
	}
	if err.Error() != "failed to load configuration after 3 attempts" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// Test: Missing JWT_SECRET produces error

func TestLoadConfig_MissingSecret(t *testing.T) {
	// Arrange
	fakeReader := FakeReader{
		Values: map[string]string{}, // no JWT_SECRET
	}

	fakeLoad := func() error {
		return nil // loading OK
	}

	loader := NewEnvLoader("ignored", fakeLoad, fakeReader)

	// Act
	cfg, err := loader.LoadConfig(0, 3)

	// Assert
	if err == nil {
		t.Fatalf("expected error because JWT_SECRET is missing, got nil")
	}
	if cfg != nil {
		t.Errorf("expected cfg to be nil, got %+v", cfg)
	}
	expected := "JWT_SECRET is missing in environment"
	if err.Error() != expected {
		t.Errorf("expected error '%s', got '%s'", expected, err.Error())
	}
}
