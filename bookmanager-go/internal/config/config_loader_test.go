package config

import (
	"errors"
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
