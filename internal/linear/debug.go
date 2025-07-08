package linear

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type DebugLogger struct {
	logger  *log.Logger
	enabled bool
}

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

func (d *DebugLogger) LogRequest(method, url, query string, variables map[string]interface{}) {
	if !d.enabled || d.logger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	d.logger.Printf("[%s] DEBUG: Linear API Request\n", timestamp)
	d.logger.Printf("  Method: %s\n", method)
	d.logger.Printf("  URL: %s\n", url)

	// Log the GraphQL query
	if query != "" {
		d.logger.Printf("  Query: %s\n", query)
	}

	if len(variables) > 0 {
		varsJSON, _ := json.MarshalIndent(variables, "  ", "  ")
		d.logger.Printf("  Variables: %s\n", string(varsJSON))
	}
	d.logger.Println()
}

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
		d.logger.Printf("  Response: %d bytes\n", len(responseBody))
	}
	d.logger.Println()
}

func (d *DebugLogger) LogError(context string, err error) {
	if !d.enabled || d.logger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	d.logger.Printf("[%s] ERROR: %s - %v", timestamp, context, err)
	d.logger.Println()
}

func (d *DebugLogger) LogInfo(format string, args ...interface{}) {
	if !d.enabled || d.logger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	d.logger.Printf("[%s] INFO: %s", timestamp, message)
	d.logger.Println()
}
