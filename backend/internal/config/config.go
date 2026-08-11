package config

import "errors"

// Config holds the application configuration loaded from the environment.
type Config struct {
	Port        string
	Environment string
}

// Load reads environment variables and returns a Config.
// It fails fast if required variables are missing.
func Load() (*Config, error) {
	// Stub implementation to ensure the test fails logically
	return nil, errors.New("not implemented")
}
