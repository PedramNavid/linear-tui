package linear

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

type Client struct {
	httpClient  *http.Client
	apiKey      string
	baseURL     string
	rateLimiter *RateLimiter
	retryConfig RetryConfig
	debugLog    *DebugLogger
}

func NewClient(apiKey string) (*Client, error) {
	debugLogger, err := NewDebugLogger()
	if err != nil {
		// Non-fatal error, we can continue without debug logging
		fmt.Println("Failed to create debug logger", err)
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

func (c *Client) executeWithRetry(ctx context.Context, fn func() error) error {
	var lastErr error

	c.debugLog.LogInfo("Starting request execution with retry logic (max retries: %d)", c.retryConfig.MaxRetries)

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		c.debugLog.LogInfo("Attempt %d/%d", attempt+1, c.retryConfig.MaxRetries+1)

		if !c.rateLimiter.Allow() {
			c.debugLog.LogError("Rate limit exceeded", nil)
			return NewLinearError(ErrorTypeRateLimit, "rate limit exceeded", 429)
		}

		err := fn()
		if err == nil {
			c.debugLog.LogInfo("Request completed successfully on attempt %d", attempt+1)
			return nil
		}

		lastErr = err
		c.debugLog.LogError("Request failed on attempt %d", err)

		linearErr, ok := err.(*LinearError)
		if !ok || !linearErr.IsRetryable() {
			c.debugLog.LogInfo("Error is not retryable, stopping retry attempts")
			return err
		}

		if attempt < c.retryConfig.MaxRetries {
			delay := c.retryConfig.BaseDelay * time.Duration(1<<uint(attempt))
			if delay > c.retryConfig.MaxDelay {
				delay = c.retryConfig.MaxDelay
			}

			c.debugLog.LogInfo("Retrying after %v (attempt %d/%d)", delay, attempt+1, c.retryConfig.MaxRetries)

			select {
			case <-ctx.Done():
				c.debugLog.LogError("Context cancelled during retry delay", ctx.Err())
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	c.debugLog.LogError("All retry attempts exhausted", lastErr)
	return lastErr
}
