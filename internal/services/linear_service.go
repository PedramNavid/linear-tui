package services

import (
	"context"
	"fmt"
	"time"

	"github.com/linear-tui/linear-tui/internal/adapters"
	"github.com/linear-tui/linear-tui/internal/config"
	"github.com/linear-tui/linear-tui/internal/domain"
	"github.com/linear-tui/linear-tui/internal/linear"
)

// LinearService handles all Linear API interactions and data conversion
type LinearService struct {
	client        *linear.Client
	adapter       *adapters.LinearAdapter
	defaultTeam   *linear.Team
	teams         []linear.Team
	users         []linear.User
	lastDataFetch time.Time
}

// NewLinearService creates a new LinearService
func NewLinearService(cfg *config.Config) (*LinearService, error) {
	if cfg.LinearAPIKey == "" {
		return nil, fmt.Errorf("linear API key not configured")
	}

	client, err := linear.NewClient(cfg.LinearAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Linear client: %w", err)
	}

	service := &LinearService{
		client:  client,
		adapter: adapters.NewLinearAdapter(),
	}

	// Initialize basic data
	if err := service.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize service: %w", err)
	}

	return service, nil
}

// initialize fetches basic workspace data needed for operations
func (s *LinearService) initialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Validate API key first
	if err := s.client.ValidateAPIKey(ctx); err != nil {
		return fmt.Errorf("API key validation failed: %w", err)
	}

	// Fetch teams
	teams, err := s.client.GetTeams(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch teams: %w", err)
	}
	s.teams = teams

	// Set default team (first available team)
	if len(s.teams) > 0 {
		s.defaultTeam = &s.teams[0]
	}

	// Fetch users for assignee lookups
	users, err := s.client.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}
	s.users = users

	s.lastDataFetch = time.Now()
	return nil
}

// GetTickets fetches issues from Linear and converts them to domain Issues for UI usage
func (s *LinearService) GetTickets() ([]domain.Issue, error) {
	if s.defaultTeam == nil {
		return nil, fmt.Errorf("no default team available - please check your Linear workspace access")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	issues, err := s.client.GetIssues(ctx, s.defaultTeam.ID, 50) // Fetch up to 50 issues
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues from Linear API: %w", err)
	}

	// Convert to domain issues for UI usage
	uiIssues := s.adapter.ConvertIssuesToUIModels(issues)
	return uiIssues, nil
}

// GetTicketByID fetches a single issue by ID from Linear and converts it to domain Issue
func (s *LinearService) GetTicketByID(issueID string) (*domain.Issue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	issue, err := s.client.GetIssueByID(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue from Linear API: %w", err)
	}

	// Convert to domain issue for UI usage
	uiIssue := s.adapter.ConvertIssueToUIModel(*issue)
	return &uiIssue, nil
}

// GetProjects fetches projects from Linear and converts them to domain Projects for UI usage
func (s *LinearService) GetProjects() ([]domain.Project, error) {
	if s.defaultTeam == nil {
		return nil, fmt.Errorf("no default team available - please check your Linear workspace access")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	projects, err := s.client.GetProjects(ctx, s.defaultTeam.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects from Linear API: %w", err)
	}

	// Convert to domain projects for UI usage
	uiProjects := s.adapter.ConvertProjectsToUIModels(projects)
	return uiProjects, nil
}

// CreateTicket creates a new issue in Linear
func (s *LinearService) CreateTicket(title, description, priority, assigneeName string) (*domain.Issue, error) {
	if s.defaultTeam == nil {
		return nil, fmt.Errorf("no default team available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert priority string to number
	priorityNum := s.adapter.ConvertPriorityToNumber(priority)

	// Find assignee ID if specified
	var assigneeID string
	if assigneeName != "" && assigneeName != "Unassigned" {
		for _, user := range s.users {
			if user.Name == assigneeName {
				assigneeID = user.ID
				break
			}
		}
	}

	input := linear.CreateIssueInput{
		Title:       title,
		Description: description,
		TeamID:      s.defaultTeam.ID,
		Priority:    priorityNum,
		AssigneeID:  assigneeID,
	}

	issue, err := s.client.CreateIssue(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	// Convert to domain issue for UI usage
	uiIssue := s.adapter.ConvertIssueToUIModel(*issue)
	return &uiIssue, nil
}

// GetTeams returns available teams
func (s *LinearService) GetTeams() []linear.Team {
	return s.teams
}

// GetUsers returns available users
func (s *LinearService) GetUsers() []linear.User {
	return s.users
}

// GetDefaultTeam returns the default team
func (s *LinearService) GetDefaultTeam() *linear.Team {
	return s.defaultTeam
}

// SetDefaultTeam sets the default team for operations
func (s *LinearService) SetDefaultTeam(teamID string) error {
	for _, team := range s.teams {
		if team.ID == teamID {
			s.defaultTeam = &team
			return nil
		}
	}
	return fmt.Errorf("team with ID %s not found", teamID)
}

// UpdateTicket updates an existing issue in Linear
func (s *LinearService) UpdateTicket(issueID, title, description, priority, assigneeName, statusName string) (*domain.Issue, error) {
	if s.defaultTeam == nil {
		return nil, fmt.Errorf("no default team available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build update input
	input := linear.UpdateIssueInput{
		Title:       title,
		Description: description,
	}

	// Convert priority string to number if provided
	if priority != "" {
		input.Priority = s.adapter.ConvertPriorityToNumber(priority)
	}

	// Find assignee ID if specified
	if assigneeName != "" && assigneeName != "Unassigned" {
		for _, user := range s.users {
			if user.Name == assigneeName {
				input.AssigneeID = user.ID
				break
			}
		}
	}

	// Find state ID if status name provided
	if statusName != "" {
		states, err := s.GetIssueStates()
		if err == nil {
			for _, state := range states {
				if state.Name == statusName {
					input.StateID = state.ID
					break
				}
			}
		}
	}

	// Update the issue
	issue, err := s.client.UpdateIssue(ctx, issueID, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue: %w", err)
	}

	// Convert to domain issue for UI usage
	uiIssue := s.adapter.ConvertIssueToUIModel(*issue)
	return &uiIssue, nil
}

// RefreshData forces a refresh of cached data
func (s *LinearService) RefreshData() error {
	return s.initialize()
}

// IsDataStale checks if cached data should be refreshed
func (s *LinearService) IsDataStale() bool {
	return time.Since(s.lastDataFetch) > 5*time.Minute
}

// GetIssueStates fetches available issue states for the default team
func (s *LinearService) GetIssueStates() ([]linear.IssueState, error) {
	if s.defaultTeam == nil {
		return nil, fmt.Errorf("no default team available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	states, err := s.client.GetIssueStates(ctx, s.defaultTeam.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue states: %w", err)
	}

	return states, nil
}
