package linear

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// Client is the main Linear API client
type Client struct {
	httpClient  *http.Client
	apiKey      string
	baseURL     string
	rateLimiter *RateLimiter
	retryConfig RetryConfig
	debugLog    *DebugLogger
}

// NewClient creates a new Linear API client
func NewClient(apiKey string) (*Client, error) {
	debugLogger, err := NewDebugLogger()
	if err != nil {
		// Non-fatal error, we can continue without debug logging
		debugLogger = &DebugLogger{enabled: false}
	}

	if apiKey == "" {
		debugLogger.LogInfo("No API key provided, checking LINEAR_API_KEY environment variable")
		// Try to get from environment variable
		apiKey = os.Getenv("LINEAR_API_KEY")
		if apiKey == "" {
			debugLogger.LogError("API key initialization", fmt.Errorf("no API key provided and LINEAR_API_KEY environment variable not set"))
			return nil, fmt.Errorf("no API key provided and LINEAR_API_KEY environment variable not set")
		}
		debugLogger.LogInfo("API key found in LINEAR_API_KEY environment variable")
	} else {
		debugLogger.LogInfo("API key provided directly to NewClient")
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:      apiKey,
		baseURL:     "https://api.linear.app/graphql",
		rateLimiter: NewRateLimiter(),
		retryConfig: RetryConfig{
			MaxRetries: 3,
			BaseDelay:  1 * time.Second,
			MaxDelay:   10 * time.Second,
		},
		debugLog: debugLogger,
	}

	debugLogger.LogInfo("Linear API client initialized")
	return client, nil
}

// executeWithRetry executes a function with retry logic
func (c *Client) executeWithRetry(ctx context.Context, fn func() error) error {
	var lastErr error

	c.debugLog.LogInfo("Starting request execution with retry logic (max retries: %d)", c.retryConfig.MaxRetries)

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		c.debugLog.LogInfo("Attempt %d/%d", attempt+1, c.retryConfig.MaxRetries+1)

		// Check rate limit
		if !c.rateLimiter.Allow() {
			c.debugLog.LogError("Rate limit exceeded", nil)
			return NewLinearError(ErrorTypeRateLimit, "rate limit exceeded", 429)
		}
		c.debugLog.LogInfo("Rate limit check passed")

		// Execute function
		err := fn()
		if err == nil {
			c.debugLog.LogInfo("Request completed successfully on attempt %d", attempt+1)
			return nil
		}

		lastErr = err
		c.debugLog.LogError("Request failed on attempt %d", err)

		// Check if error is retryable
		linearErr, ok := err.(*LinearError)
		if !ok || !linearErr.IsRetryable() {
			c.debugLog.LogInfo("Error is not retryable, stopping retry attempts")
			return err
		}

		if attempt < c.retryConfig.MaxRetries {
			// Calculate delay with exponential backoff
			delay := c.retryConfig.BaseDelay * time.Duration(1<<uint(attempt))
			if delay > c.retryConfig.MaxDelay {
				delay = c.retryConfig.MaxDelay
			}

			c.debugLog.LogInfo("Retrying after %v (attempt %d/%d)", delay, attempt+1, c.retryConfig.MaxRetries)

			// Check context cancellation
			select {
			case <-ctx.Done():
				c.debugLog.LogError("Context cancelled during retry delay", ctx.Err())
				return ctx.Err()
			case <-time.After(delay):
				// Continue to next attempt
			}
		}
	}

	c.debugLog.LogError("All retry attempts exhausted", lastErr)
	return lastErr
}
