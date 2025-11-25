package config

// ConfigLoader defines an interface for loading configuration.
type ConfigLoader interface {
	LoadConfig(attempts, maxAttempts int) error
}
