package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	Environment string
}

func Load() (*Config, error) {
	port, exists := os.LookupEnv("PORT")
	if !exists || port == "" {
		return nil, errors.New("missing required environment variable: PORT")
	}

	env, exists := os.LookupEnv("ENV")
	if !exists || env == "" {
		// Defaulting to development if not explicitly set,
		env = "development"
	}

	return &Config{
		Port:        port,
		Environment: env,
	}, nil
}
