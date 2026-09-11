package errors_test

import (
	"encoding/json"
	"testing"

	"github.com/vrushaliraut/insta-gate/backend/pkg/errors"
)

func TestAPIErrorJSONSerialization(t *testing.T) {
	// Create a new Structured error
	apiErr := errors.NewAPIError("UNAUTHORIZED", "Invalid or missing token")

	// Marshal it to JSON
	bytes, err := json.Marshal(apiErr)
	if err != nil {
		t.Fatalf("Failed to marshal APIError: %v", err)
	}

	// Unmarshal back into a generic map to verify the exact structure
	var result map[string]map[string]string
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON into map: %v", err)
	}

	// Assert the nested structured matches: {"error": {"code": "", "message": ""}}
	errObj, exists := result["error"]
	if !exists {
		t.Fatalf("Expected top-level 'error' key, but it was missing")
	}

	if errObj["code"] != "UNAUTHORIZED" {
		t.Errorf("Expected code 'UNAUTHORIZED', got '%s'", errObj["code"])
	}

	if errObj["message"] != "Invalid or missing token" {
		t.Errorf("Expected message 'Invalid or missing token', got '%s'", errObj["message"])
	}
}
