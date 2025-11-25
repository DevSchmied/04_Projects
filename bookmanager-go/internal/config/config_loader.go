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

// EnvLoader loads environment variables using a custom load function.
type EnvLoader struct {
	FilePath string
	Load     func() error
}

// NewEnvLoader initializes a new EnvLoader.
func NewEnvLoader(path string, loadFunc func() error) *EnvLoader {
	return &EnvLoader{
		FilePath: path,
		Load:     loadFunc,
	}
}

// LoadConfig loads configuration recursively with a retry mechanism.
func (e *EnvLoader) LoadConfig(attempt, maxAttempts int) error {
	if attempt >= maxAttempts {
		return fmt.Errorf("failed to load configuration after %d attempts", maxAttempts)
	}

	// Try to load .env
	err := e.Load()
	if err == nil {
		// Read required fields
		Cfg.JWTSecret = os.Getenv("JWT_SECRET")

		if Cfg.JWTSecret == "" {
			return fmt.Errorf("JWT_SECRET is missing in environment")
		}

		return nil
	}

	log.Printf("Error loading config (attempt %d): %v\n", attempt+1, err)

	// Retry mechanism
	return e.LoadConfig(attempt, maxAttempts)
}
