package adapters

import (
	"time"
	"github.com/linear-tui/linear-tui/internal/linear"
	"github.com/linear-tui/linear-tui/internal/ui/mock"
)

// LinearAdapter converts between Linear API types and UI types
type LinearAdapter struct{}

// NewLinearAdapter creates a new LinearAdapter
func NewLinearAdapter() *LinearAdapter {
	return &LinearAdapter{}
}

// ConvertIssueToMockTicket converts a Linear Issue to a MockTicket for UI compatibility
func (a *LinearAdapter) ConvertIssueToMockTicket(issue linear.Issue) mock.MockTicket {
	// Convert priority number to string
	priorityStr := a.convertPriorityToString(issue.Priority)

	// Get assignee name or default
	assigneeName := "Unassigned"
	if issue.Assignee != nil {
		assigneeName = issue.Assignee.Name
	}

	return mock.MockTicket{
		ID:          issue.ID,
		Title:       issue.Title,
		Description: issue.Description,
		Status:      issue.State.Name,
		Priority:    priorityStr,
		Assignee:    assigneeName,
		CreatedAt:   issue.CreatedAt,
	}
}

// ConvertIssuesToMockTickets converts a slice of Linear Issues to MockTickets
func (a *LinearAdapter) ConvertIssuesToMockTickets(issues []linear.Issue) []mock.MockTicket {
	tickets := make([]mock.MockTicket, len(issues))
	for i, issue := range issues {
		tickets[i] = a.ConvertIssueToMockTicket(issue)
	}
	return tickets
}

// ConvertProjectToMockProject converts a Linear Project to a MockProject for UI compatibility
func (a *LinearAdapter) ConvertProjectToMockProject(project linear.Project) mock.MockProject {
	return mock.MockProject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Status:      project.State,
		Progress:    project.Progress,
		CreatedAt:   parseProjectDate(project.StartDate), // Use start date as created date
	}
}

// ConvertProjectsToMockProjects converts a slice of Linear Projects to MockProjects
func (a *LinearAdapter) ConvertProjectsToMockProjects(projects []linear.Project) []mock.MockProject {
	mockProjects := make([]mock.MockProject, len(projects))
	for i, project := range projects {
		mockProjects[i] = a.ConvertProjectToMockProject(project)
	}
	return mockProjects
}

// ConvertMockTicketToCreateIssueInput converts a MockTicket to CreateIssueInput
func (a *LinearAdapter) ConvertMockTicketToCreateIssueInput(ticket mock.MockTicket, teamID string) linear.CreateIssueInput {
	return linear.CreateIssueInput{
		Title:       ticket.Title,
		Description: ticket.Description,
		TeamID:      teamID,
		Priority:    a.ConvertPriorityToNumber(ticket.Priority),
		// Note: StateID and AssigneeID would need to be resolved separately
		// since the mock types don't contain the actual IDs
	}
}

// convertPriorityToString converts Linear priority number to string
func (a *LinearAdapter) convertPriorityToString(priority int) string {
	switch priority {
	case 0:
		return "None"
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Normal"
	case 4:
		return "Low"
	default:
		return "Unknown"
	}
}

// ConvertPriorityToNumber converts priority string to Linear priority number
func (a *LinearAdapter) ConvertPriorityToNumber(priority string) int {
	switch priority {
	case "None":
		return 0
	case "Urgent":
		return 1
	case "High":
		return 2
	case "Normal", "Medium":
		return 3
	case "Low":
		return 4
	default:
		return 0
	}
}

// GetDefaultTeamID returns a default team ID for operations that require one
// This is a placeholder - in a real implementation, this would be configurable
func (a *LinearAdapter) GetDefaultTeamID() string {
	// This should be configurable or determined from the user's workspace
	return ""
}

// parseProjectDate converts a Linear project date string to time.Time
// Linear returns date strings like "2025-07-06" which need special parsing
func parseProjectDate(dateStr *string) time.Time {
	if dateStr == nil || *dateStr == "" {
		return time.Time{} // Return zero time for nil/empty dates
	}
	
	// Parse date-only format "2006-01-02"
	parsed, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return time.Time{} // Return zero time if parsing fails
	}
	
	return parsed
}
