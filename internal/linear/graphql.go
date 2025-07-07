package linear

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message    string                 `json:"message"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Path       []interface{}          `json:"path,omitempty"`
}

// executeGraphQL executes a GraphQL query/mutation
func (c *Client) executeGraphQL(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	// Create request body
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)

	// Log request in debug mode
	c.debugLog.LogRequest("POST", c.baseURL, query, variables)

	// Execute request
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		c.debugLog.LogResponse(0, duration, nil, err)
		return NewLinearError(ErrorTypeNetwork, fmt.Sprintf("network error: %v", err), 0)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.debugLog.LogResponse(resp.StatusCode, duration, nil, fmt.Errorf("failed to close response body: %w", err))
		}
	}()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.debugLog.LogResponse(resp.StatusCode, duration, nil, err)
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log response in debug mode
	c.debugLog.LogResponse(resp.StatusCode, duration, body, nil)

	// Check HTTP status code
	if resp.StatusCode == 401 {
		return NewLinearError(ErrorTypeAuth, "authentication failed - invalid API key", 401)
	}

	if resp.StatusCode == 429 {
		return NewLinearError(ErrorTypeRateLimit, "rate limit exceeded", 429)
	}

	if resp.StatusCode != 200 {
		return NewLinearError(ErrorTypeAPI, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), resp.StatusCode)
	}

	// Parse GraphQL response
	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for GraphQL errors
	if len(gqlResp.Errors) > 0 {
		// Combine all error messages
		var messages []string
		for _, e := range gqlResp.Errors {
			messages = append(messages, e.Message)
		}
		return NewLinearError(ErrorTypeAPI, fmt.Sprintf("GraphQL errors: %v", messages), 200)
	}

	// Parse data into result
	if result != nil && gqlResp.Data != nil {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("failed to parse data: %w", err)
		}
	}

	return nil
}

// Query executes a GraphQL query
func (c *Client) Query(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	var result json.RawMessage
	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Mutation executes a GraphQL mutation
func (c *Client) Mutation(ctx context.Context, mutation string, variables map[string]interface{}) ([]byte, error) {
	var result json.RawMessage
	if err := c.executeGraphQL(ctx, mutation, variables, &result); err != nil {
		return nil, err
	}
	return result, nil
}
