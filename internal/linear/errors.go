package linear

import (
	"fmt"
)

// ErrorType represents the type of error from the Linear API
type ErrorType string

const (
	ErrorTypeNetwork    ErrorType = "network"
	ErrorTypeAuth       ErrorType = "auth"
	ErrorTypeAPI        ErrorType = "api"
	ErrorTypeRateLimit  ErrorType = "ratelimit"
	ErrorTypeValidation ErrorType = "validation"
)

// LinearError represents an error from the Linear API
type LinearError struct {
	Type    ErrorType
	Message string
	Code    int
}

// Error implements the error interface
func (e *LinearError) Error() string {
	return fmt.Sprintf("Linear API error [%s]: %s (code: %d)", e.Type, e.Message, e.Code)
}

// IsRetryable returns true if the error is retryable
func (e *LinearError) IsRetryable() bool {
	switch e.Type {
	case ErrorTypeNetwork, ErrorTypeRateLimit:
		return true
	case ErrorTypeAPI:
		// Some API errors might be retryable (e.g., 500, 502, 503)
		return e.Code >= 500 && e.Code < 600
	default:
		return false
	}
}

// NewLinearError creates a new LinearError
func NewLinearError(errorType ErrorType, message string, code int) *LinearError {
	return &LinearError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}
