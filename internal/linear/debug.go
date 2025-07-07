package linear

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// DebugLogger handles debug logging for the Linear client
type DebugLogger struct {
	logger  *log.Logger
	enabled bool
}

// NewDebugLogger creates a new debug logger
func NewDebugLogger() (*DebugLogger, error) {
	debugLogger := &DebugLogger{
		enabled: len(os.Getenv("DEBUG")) > 0,
	}

	if debugLogger.enabled {
		file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open debug.log: %w", err)
		}

		debugLogger.logger = log.New(file, "", 0)
	}

	return debugLogger, nil
}

// LogRequest logs an API request
func (d *DebugLogger) LogRequest(method, url, query string, variables map[string]interface{}) {
	if !d.enabled || d.logger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	d.logger.Printf("[%s] DEBUG: Linear API Request\n", timestamp)
	d.logger.Printf("  Method: %s\n", method)
	d.logger.Printf("  URL: %s\n", url)
	// Query removed from logging for security/verbosity reasons

	if len(variables) > 0 {
		varsJSON, _ := json.MarshalIndent(variables, "  ", "  ")
		d.logger.Printf("  Variables: %s\n", string(varsJSON))
	}
	d.logger.Println()
}

// LogResponse logs an API response
func (d *DebugLogger) LogResponse(statusCode int, duration time.Duration, responseBody []byte, err error) {
	if !d.enabled || d.logger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	d.logger.Printf("[%s] DEBUG: Linear API Response\n", timestamp)
	d.logger.Printf("  Status: %d\n", statusCode)
	d.logger.Printf("  Duration: %v\n", duration)

	if err != nil {
		d.logger.Printf("  Error: %v\n", err)
	} else if len(responseBody) > 0 {
		// Just log response size instead of full body
		d.logger.Printf("  Response: %d bytes\n", len(responseBody))
	}
	d.logger.Println()
}

// LogError logs an error
func (d *DebugLogger) LogError(context string, err error) {
	if !d.enabled || d.logger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	d.logger.Printf("[%s] ERROR: %s - %v", timestamp, context, err)
	d.logger.Println()
}

// LogInfo logs general information
func (d *DebugLogger) LogInfo(format string, args ...interface{}) {
	if !d.enabled || d.logger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	d.logger.Printf("[%s] INFO: %s", timestamp, message)
	d.logger.Println()
}
