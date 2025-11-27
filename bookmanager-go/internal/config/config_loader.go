package config

import (
	"fmt"
	"log"
	"os"
)

// ConfigLoader defines an interface for loading configuration.
type ConfigLoader interface {
	LoadConfig(attempt, maxAttempts int) error
}

// EnvLoader loads configuration via injected Load function & environment reader.
type EnvLoader struct {
	FilePath string
	Load     func() error
	Reader   EnvReader
}

// EnvReader abstracts environment variable access for testability.
type EnvReader interface {
	Get(key string) string
}

// OSReader is the default implementation using os.Getenv.
type OSReader struct{}

// Get retrieves an environment variable using os.Getenv.
func (OSReader) Get(key string) string {
	return os.Getenv(key)
}

// NewEnvLoader initializes a new EnvLoader.
func NewEnvLoader(path string, loadFunc func() error, reader EnvReader) *EnvLoader {
	return &EnvLoader{
		FilePath: path,
		Load:     loadFunc,
		Reader:   reader,
	}
}

// LoadConfig loads configuration recursively with a retry mechanism and returns AppConfig.
func (e *EnvLoader) LoadConfig(attempt, maxAttempts int) (*AppConfig, error) {
	if attempt >= maxAttempts {
		return nil, fmt.Errorf("failed to load configuration after %d attempts", maxAttempts)
	}

	// Try to load .env
	err := e.Load()
	if err == nil {
		// Read required fields
		cfg := &AppConfig{
			JWTSecret: e.Reader.Get("JWT_SECRET"),
		}

		if cfg.JWTSecret == "" {
			return nil, fmt.Errorf("JWT_SECRET is missing in environment")
		}

		return cfg, nil
	}

	log.Printf("Error loading config (attempt %d): %v\n", attempt+1, err)

	// Retry mechanism
	return e.LoadConfig(attempt+1, maxAttempts)
}
