package linear

import (
	"context"
	"fmt"
)

// GetIssues retrieves issues for a team
func (c *Client) GetIssues(ctx context.Context, teamID string, limit int) ([]Issue, error) {
	c.debugLog.LogInfo("Fetching issues for team %s (limit: %d)", teamID, limit)
	query := `
		query GetIssues($teamId: ID!, $first: Int!) {
			issues(filter: { team: { id: { eq: $teamId } } }, first: $first) {
				nodes {
					id
					title
					description
					priority
					createdAt
					updatedAt
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
						avatarUrl
					}
					team {
						id
						name
						key
					}
					project {
						id
						name
					}
				}
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"teamId": teamID,
		"first":  limit,
	}

	var response IssuesResponse
	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, query, variables, &response)
	})

	if err != nil {
		c.debugLog.LogError("Failed to fetch issues", err)
		return nil, err
	}

	c.debugLog.LogInfo("Successfully fetched %d issues for team %s", len(response.Issues.Nodes), teamID)
	return response.Issues.Nodes, nil
}

// GetProjects retrieves projects for a team
func (c *Client) GetProjects(ctx context.Context, teamID string) ([]Project, error) {
	c.debugLog.LogInfo("Fetching projects for team %s", teamID)
	query := `
		query GetProjects {
			projects {
				nodes {
					id
					name
					description
					state
					progress
					startDate
					targetDate
				}
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}
	`

	var response ProjectsResponse
	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, query, nil, &response)
	})

	if err != nil {
		c.debugLog.LogError("Failed to fetch projects", err)
		return nil, err
	}

	c.debugLog.LogInfo("Successfully fetched %d projects for team %s", len(response.Projects.Nodes), teamID)
	return response.Projects.Nodes, nil
}

// GetTeams retrieves all teams
func (c *Client) GetTeams(ctx context.Context) ([]Team, error) {
	c.debugLog.LogInfo("Fetching all teams")
	query := `
		query GetTeams {
			teams {
				nodes {
					id
					name
					description
					key
				}
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}
	`

	var response TeamsResponse
	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, query, nil, &response)
	})

	if err != nil {
		c.debugLog.LogError("Failed to fetch teams", err)
		return nil, err
	}

	c.debugLog.LogInfo("Successfully fetched %d teams", len(response.Teams.Nodes))
	return response.Teams.Nodes, nil
}

// GetUsers retrieves all users
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	query := `
		query GetUsers {
			users {
				nodes {
					id
					name
					email
					avatarUrl
				}
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}
	`

	var response UsersResponse
	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, query, nil, &response)
	})

	if err != nil {
		return nil, err
	}

	return response.Users.Nodes, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, input CreateIssueInput) (*Issue, error) {
	mutation := `
		mutation CreateIssue($input: IssueCreateInput!) {
			issueCreate(input: $input) {
				success
				issue {
					id
					title
					description
					priority
					createdAt
					updatedAt
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
						avatarUrl
					}
					team {
						id
						name
						key
					}
					project {
						id
						name
					}
				}
			}
		}
	`

	// Build input variables
	inputMap := map[string]interface{}{
		"title":  input.Title,
		"teamId": input.TeamID,
	}

	if input.Description != "" {
		inputMap["description"] = input.Description
	}
	if input.Priority > 0 {
		inputMap["priority"] = input.Priority
	}
	if input.AssigneeID != "" {
		inputMap["assigneeId"] = input.AssigneeID
	}
	if input.ProjectID != "" {
		inputMap["projectId"] = input.ProjectID
	}
	if input.StateID != "" {
		inputMap["stateId"] = input.StateID
	}

	variables := map[string]interface{}{
		"input": inputMap,
	}

	var response struct {
		IssueCreate struct {
			Success bool  `json:"success"`
			Issue   Issue `json:"issue"`
		} `json:"issueCreate"`
	}

	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, mutation, variables, &response)
	})

	if err != nil {
		return nil, err
	}

	if !response.IssueCreate.Success {
		return nil, NewLinearError(ErrorTypeAPI, "failed to create issue", 200)
	}

	return &response.IssueCreate.Issue, nil
}

// UpdateIssue updates an existing issue
func (c *Client) UpdateIssue(ctx context.Context, id string, input UpdateIssueInput) (*Issue, error) {
	mutation := `
		mutation UpdateIssue($id: ID!, $input: IssueUpdateInput!) {
			issueUpdate(id: $id, input: $input) {
				success
				issue {
					id
					title
					description
					priority
					createdAt
					updatedAt
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
						avatarUrl
					}
					team {
						id
						name
						key
					}
					project {
						id
						name
					}
				}
			}
		}
	`

	// Build input variables
	inputMap := make(map[string]interface{})

	if input.Title != "" {
		inputMap["title"] = input.Title
	}
	if input.Description != "" {
		inputMap["description"] = input.Description
	}
	if input.Priority > 0 {
		inputMap["priority"] = input.Priority
	}
	if input.AssigneeID != "" {
		inputMap["assigneeId"] = input.AssigneeID
	}
	if input.ProjectID != "" {
		inputMap["projectId"] = input.ProjectID
	}
	if input.StateID != "" {
		inputMap["stateId"] = input.StateID
	}

	variables := map[string]interface{}{
		"id":    id,
		"input": inputMap,
	}

	var response struct {
		IssueUpdate struct {
			Success bool  `json:"success"`
			Issue   Issue `json:"issue"`
		} `json:"issueUpdate"`
	}

	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, mutation, variables, &response)
	})

	if err != nil {
		return nil, err
	}

	if !response.IssueUpdate.Success {
		return nil, NewLinearError(ErrorTypeAPI, "failed to update issue", 200)
	}

	return &response.IssueUpdate.Issue, nil
}

// CreateComment creates a comment on an issue
func (c *Client) CreateComment(ctx context.Context, issueID, body string) (*Comment, error) {
	mutation := `
		mutation CreateComment($input: CommentCreateInput!) {
			commentCreate(input: $input) {
				success
				comment {
					id
					body
					createdAt
					updatedAt
					user {
						id
						name
						email
						avatarUrl
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"issueId": issueID,
			"body":    body,
		},
	}

	var response struct {
		CommentCreate struct {
			Success bool    `json:"success"`
			Comment Comment `json:"comment"`
		} `json:"commentCreate"`
	}

	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, mutation, variables, &response)
	})

	if err != nil {
		return nil, err
	}

	if !response.CommentCreate.Success {
		return nil, NewLinearError(ErrorTypeAPI, "failed to create comment", 200)
	}

	return &response.CommentCreate.Comment, nil
}

// GetIssueStates retrieves available issue states for a team
func (c *Client) GetIssueStates(ctx context.Context, teamID string) ([]IssueState, error) {
	query := `
		query GetIssueStates($teamId: ID!) {
			team(id: $teamId) {
				states {
					nodes {
						id
						name
						type
						color
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"teamId": teamID,
	}

	var response struct {
		Team struct {
			States struct {
				Nodes []IssueState `json:"nodes"`
			} `json:"states"`
		} `json:"team"`
	}

	err := c.executeWithRetry(ctx, func() error {
		return c.executeGraphQL(ctx, query, variables, &response)
	})

	if err != nil {
		return nil, err
	}

	return response.Team.States.Nodes, nil
}

// ValidateAPIKey validates the API key by making a simple query
func (c *Client) ValidateAPIKey(ctx context.Context) error {
	c.debugLog.LogInfo("Starting API key validation")

	query := `
		query ValidateAPIKey {
			viewer {
				id
				name
				email
			}
		}
	`

	var response struct {
		Viewer struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"viewer"`
	}

	err := c.executeGraphQL(ctx, query, nil, &response)
	if err != nil {
		c.debugLog.LogError("API key validation", err)
		return fmt.Errorf("API key validation failed: %w", err)
	}

	c.debugLog.LogInfo("API key validated successfully for user: %s (%s) [ID: %s]",
		response.Viewer.Name, response.Viewer.Email, response.Viewer.ID)
	return nil
}
