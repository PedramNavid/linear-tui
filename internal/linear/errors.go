package linear

import (
	"fmt"
)

type ErrorType string

const (
	ErrorTypeNetwork    ErrorType = "network"
	ErrorTypeAuth       ErrorType = "auth"
	ErrorTypeAPI        ErrorType = "api"
	ErrorTypeRateLimit  ErrorType = "ratelimit"
	ErrorTypeValidation ErrorType = "validation"
)

type LinearError struct {
	Type    ErrorType
	Message string
	Code    int
}

func (e *LinearError) Error() string {
	return fmt.Sprintf("Linear API error [%s]: %s (code: %d)", e.Type, e.Message, e.Code)
}

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

func NewLinearError(errorType ErrorType, message string, code int) *LinearError {
	return &LinearError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}
