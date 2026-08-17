package errors

// ErrorDetail holds the specific error information.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIError is the standard JSON error response wrapper
type APIError struct {
	Error ErrorDetail `json:"error"`
}

// NewAPIError creates a properly formatted API error response.
func NewAPIError(code, message string) *APIError {
	return &APIError{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
