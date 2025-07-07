package services

import (
	"context"
	"fmt"
	"time"

	"github.com/linear-tui/linear-tui/internal/config"
	"github.com/linear-tui/linear-tui/internal/linear"
	"github.com/linear-tui/linear-tui/internal/ui/adapters"
	"github.com/linear-tui/linear-tui/internal/ui/mock"
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

// GetTickets fetches issues from Linear and converts them to MockTickets for UI compatibility
func (s *LinearService) GetTickets() ([]mock.MockTicket, error) {
	if s.defaultTeam == nil {
		return nil, fmt.Errorf("no default team available - please check your Linear workspace access")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	issues, err := s.client.GetIssues(ctx, s.defaultTeam.ID, 50) // Fetch up to 50 issues
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues from Linear API: %w", err)
	}

	// Convert to mock tickets for UI compatibility
	tickets := s.adapter.ConvertIssuesToMockTickets(issues)
	return tickets, nil
}

// GetProjects fetches projects from Linear and converts them to MockProjects for UI compatibility
func (s *LinearService) GetProjects() ([]mock.MockProject, error) {
	if s.defaultTeam == nil {
		return nil, fmt.Errorf("no default team available - please check your Linear workspace access")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	projects, err := s.client.GetProjects(ctx, s.defaultTeam.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects from Linear API: %w", err)
	}

	// Convert to mock projects for UI compatibility
	mockProjects := s.adapter.ConvertProjectsToMockProjects(projects)
	return mockProjects, nil
}

// CreateTicket creates a new issue in Linear
func (s *LinearService) CreateTicket(title, description, priority, assigneeName string) (*mock.MockTicket, error) {
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

	// Convert to mock ticket for UI compatibility
	ticket := s.adapter.ConvertIssueToMockTicket(*issue)
	return &ticket, nil
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
