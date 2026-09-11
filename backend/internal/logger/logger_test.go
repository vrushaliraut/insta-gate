package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/vrushaliraut/insta-gate/backend/internal/logger"
)

func TestLoggerJSONFormatAndRequestID(t *testing.T) {
	// Setup a buffer to capture log output instead of printing to stdout
	var buf bytes.Buffer

	// Initialize the logger
	log := logger.New(&buf, "development")

	// Create a context containing mocked request ID
	ctx := context.WithValue(context.Background(), logger.RequestIDKey, "req-123-abc")

	// Log test message passing the context
	log.InfoContext(ctx, "user logged in", "user_id", 42)

	// Parse the captured output as JSON
	var logOutput map[string]interface{}

	if err := json.Unmarshal(buf.Bytes(), &logOutput); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v. Raw Output: %s", err, buf.String())
	}

	// Assertions
	if logOutput["msg"] != "user logged in" {
		t.Errorf("expected msg 'user logged in', got '%v'", logOutput["msg"])
	}

	if logOutput["user_id"] != float64(42) { // JSON unmarshals numbers to float64
		t.Errorf("expected user_id 42, got '%v'", logOutput["user_id"])
	}

	if logOutput["request_id"] != "req-123-abc" {
		t.Errorf("expected msg 'req-123-abc', got '%v'", logOutput["request_id"])
	}
}
