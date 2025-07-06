package linear

import (
	"context"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		baseURL:    "https://api.linear.app/graphql",
	}
}

func (c *Client) Query(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	// TODO: Implement GraphQL query logic
	return nil, nil
}

type Issue struct {
	ID          string
	Title       string
	Description string
	State       string
	Priority    int
	Assignee    *User
	CreatedAt   string
	UpdatedAt   string
}

type Project struct {
	ID          string
	Name        string
	Description string
	State       string
	Progress    float64
	StartDate   string
	TargetDate  string
}

type Team struct {
	ID          string
	Name        string
	Description string
	Key         string
}

type User struct {
	ID       string
	Name     string
	Email    string
	AvatarURL string
}