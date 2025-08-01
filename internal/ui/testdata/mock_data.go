package testdata

import (
	"time"

	"github.com/charmbracelet/bubbles/list"

	"github.com/linear-tui/linear-tui/internal/bubble"
	"github.com/linear-tui/linear-tui/internal/domain"
)

// GetTestIssues returns a list of test issues for development/testing
func GetTestIssues() []list.Item {

	items := []list.Item{
		bubble.Item{
			ID:          "LIN-001",
			Title:       "Implement user authentication",
			Description: "Add OAuth2 authentication flow for user login. This should include Google and GitHub providers with proper session management and token refresh capabilities.",
			Status:      "In Progress",
			Priority:    "High",
			Assignee:    "Alice Johnson",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
		bubble.Item{
			ID:          "LIN-002",
			Title:       "Fix database connection pooling",
			Description: "The current database connection pool is not properly handling timeouts and connection limits. This causes the application to hang under high load.",
			Status:      "Todo",
			Priority:    "High",
			Assignee:    "Bob Smith",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		bubble.Item{
			ID:          "LIN-003",
			Title:       "Add dark mode support",
			Description: "Implement dark mode theme throughout the application with proper color scheme and user preference storage.",
			Status:      "Done",
			Priority:    "Medium",
			Assignee:    "Carol Davis",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		},
		bubble.Item{
			ID:          "LIN-004",
			Title:       "Optimize API response times",
			Description: "Current API response times are averaging 500ms. Need to implement caching strategies and query optimization to bring this down to under 200ms.",
			Status:      "In Progress",
			Priority:    "High",
			Assignee:    "David Wilson",
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
		bubble.Item{
			ID:          "LIN-005",
			Title:       "Create user onboarding flow",
			Description: "Design and implement a comprehensive onboarding experience for new users including tutorials and tooltips.",
			Status:      "Todo",
			Priority:    "Medium",
			Assignee:    "Unassigned",
			CreatedAt:   time.Now().Add(-96 * time.Hour),
		},
	}

	return items
}

// GetTestProjects returns a list of test projects for development/testing
func GetTestProjects() []domain.Project {
	return []domain.Project{
		{
			ID:          "PRJ-001",
			Name:        "Q4 2024 Platform Improvements",
			Description: "Major platform improvements including performance optimization, security enhancements, and UI/UX updates.",
			Status:      "In Progress",
			Progress:    0.65,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		},
		{
			ID:          "PRJ-002",
			Name:        "Mobile App Development",
			Description: "Native iOS and Android applications with full feature parity with web platform.",
			Status:      "Planning",
			Progress:    0.15,
			CreatedAt:   time.Now().Add(-15 * 24 * time.Hour),
		},
		{
			ID:          "PRJ-003",
			Name:        "API v2.0",
			Description: "Complete API redesign with GraphQL support, improved documentation, and better rate limiting.",
			Status:      "In Progress",
			Progress:    0.40,
			CreatedAt:   time.Now().Add(-45 * 24 * time.Hour),
		},
	}
}
