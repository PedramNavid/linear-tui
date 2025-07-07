package mock

import "time"

// MockTicket represents a mock ticket for testing
type MockTicket struct {
	ID          string
	Title       string
	Description string
	Status      string
	Priority    string
	Assignee    string
	CreatedAt   time.Time
}

// MockProject represents a mock project for testing
type MockProject struct {
	ID          string
	Name        string
	Description string
	Status      string
	Progress    float64
	CreatedAt   time.Time
}

// GetMockTickets returns a list of mock tickets
func GetMockTickets() []MockTicket {
	return []MockTicket{
		{
			ID:          "LIN-001",
			Title:       "Implement user authentication",
			Description: "Add OAuth2 authentication flow for user login. This should include Google and GitHub providers with proper session management and token refresh capabilities.",
			Status:      "In Progress",
			Priority:    "High",
			Assignee:    "Alice Johnson",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
		{
			ID:          "LIN-002",
			Title:       "Fix database connection pooling",
			Description: "The current database connection pool is not properly handling timeouts and connection limits. This causes the application to hang under high load.",
			Status:      "Todo",
			Priority:    "High",
			Assignee:    "Bob Smith",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "LIN-003",
			Title:       "Add dark mode support",
			Description: "Implement dark mode theme throughout the application with proper color scheme and user preference storage.",
			Status:      "Done",
			Priority:    "Medium",
			Assignee:    "Carol Davis",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		},
		{
			ID:          "LIN-004",
			Title:       "Optimize search performance",
			Description: "The search functionality is slow when dealing with large datasets. Implement proper indexing and caching strategies.",
			Status:      "Backlog",
			Priority:    "Medium",
			Assignee:    "David Wilson",
			CreatedAt:   time.Now().Add(-96 * time.Hour),
		},
		{
			ID:          "LIN-005",
			Title:       "Add keyboard shortcuts",
			Description: "Users have requested keyboard shortcuts for common actions. Priority shortcuts: Ctrl+N for new item, Ctrl+S for save, Ctrl+F for search.",
			Status:      "Todo",
			Priority:    "Low",
			Assignee:    "Eve Brown",
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
	}
}

// GetMockProjects returns a list of mock projects
func GetMockProjects() []MockProject {
	return []MockProject{
		{
			ID:          "PROJ-001",
			Name:        "Authentication System",
			Description: "Complete overhaul of the authentication system with modern security practices, OAuth2 integration, and multi-factor authentication support.",
			Status:      "Active",
			Progress:    0.75,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		},
		{
			ID:          "PROJ-002",
			Name:        "Mobile App Development",
			Description: "Native mobile applications for iOS and Android with feature parity to the web application and offline capabilities.",
			Status:      "Active",
			Progress:    0.45,
			CreatedAt:   time.Now().Add(-60 * 24 * time.Hour),
		},
		{
			ID:          "PROJ-003",
			Name:        "API Documentation",
			Description: "Comprehensive API documentation with interactive examples, SDKs for popular languages, and developer onboarding guides.",
			Status:      "Completed",
			Progress:    1.0,
			CreatedAt:   time.Now().Add(-120 * 24 * time.Hour),
		},
		{
			ID:          "PROJ-004",
			Name:        "Performance Optimization",
			Description: "System-wide performance improvements including database optimization, caching strategies, and frontend bundle optimization.",
			Status:      "Planning",
			Progress:    0.0,
			CreatedAt:   time.Now().Add(-7 * 24 * time.Hour),
		},
	}
}
