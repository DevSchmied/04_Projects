package config

// ConfigLoader defines an interface for loading configuration.
type ConfigLoader interface {
	LoadConfig(attempts, maxAttempts int) error
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
