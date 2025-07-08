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

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
}

type GraphQLError struct {
	Message    string                 `json:"message"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Path       []interface{}          `json:"path,omitempty"`
}

func (c *Client) executeGraphQL(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)

	c.debugLog.LogRequest("POST", c.baseURL, query, variables)

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

	c.debugLog.LogResponse(resp.StatusCode, duration, body, nil)

	if resp.StatusCode == 401 {
		return NewLinearError(ErrorTypeAuth, "authentication failed - invalid API key", 401)
	}

	if resp.StatusCode == 429 {
		return NewLinearError(ErrorTypeRateLimit, "rate limit exceeded", 429)
	}

	if resp.StatusCode != 200 {
		// Try to parse error response
		var errorResp struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && len(errorResp.Errors) > 0 {
			return NewLinearError(ErrorTypeAPI, fmt.Sprintf("Linear API error: %s (status: %d)", errorResp.Errors[0].Message, resp.StatusCode), resp.StatusCode)
		}
		// If can't parse, return generic error with body for debugging
		return NewLinearError(ErrorTypeAPI, fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, string(body)), resp.StatusCode)
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		var messages []string
		for _, e := range gqlResp.Errors {
			messages = append(messages, e.Message)
		}
		return NewLinearError(ErrorTypeAPI, fmt.Sprintf("GraphQL errors: %v", messages), 200)
	}

	if result != nil && gqlResp.Data != nil {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("failed to parse data: %w", err)
		}
	}

	return nil
}

func (c *Client) Query(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	var result json.RawMessage
	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Mutation(ctx context.Context, mutation string, variables map[string]interface{}) ([]byte, error) {
	var result json.RawMessage
	if err := c.executeGraphQL(ctx, mutation, variables, &result); err != nil {
		return nil, err
	}
	return result, nil
}
