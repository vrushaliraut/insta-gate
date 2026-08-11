package config_test

import (
	"testing"

	"github.com/vrushaliraut/insta-gate/backend/internal/config"
)

func TestLoad_Success(t *testing.T) {
	// Setup valid environment
	t.Setenv("PORT", "8080")
	t.Setenv("ENV", "development")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected port '8080', got '%s'", cfg.Port)
	}
	if cfg.Environment != "development" {
		t.Errorf("expected environment 'development', got '%s'", cfg.Environment)
	}

}

func TestLoad_MissingRequiredVars(t *testing.T) {
	// Setup invalid environment (missing PORT)
	t.Setenv("PORT", "")
	t.Setenv("ENV", "development")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error due to missing required environment variable, got nil")
	}
}
