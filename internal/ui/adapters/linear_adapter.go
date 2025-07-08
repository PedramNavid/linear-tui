package adapters

import (
	"github.com/linear-tui/linear-tui/internal/domain"
	"github.com/linear-tui/linear-tui/internal/linear"
	"time"
)

// LinearAdapter converts between Linear API types and UI types
type LinearAdapter struct{}

// NewLinearAdapter creates a new LinearAdapter
func NewLinearAdapter() *LinearAdapter {
	return &LinearAdapter{}
}

// ConvertIssueToUIModel converts a Linear Issue to a domain Issue for UI usage
func (a *LinearAdapter) ConvertIssueToUIModel(issue linear.Issue) domain.Issue {
	// Convert priority number to string
	priorityStr := a.convertPriorityToString(issue.Priority)

	// Get assignee name or default
	assigneeName := "Unassigned"
	if issue.Assignee != nil {
		assigneeName = issue.Assignee.Name
	}

	return domain.Issue{
		ID:          issue.Identifier,
		LinearID:    issue.ID,
		Title:       issue.Title,
		Description: issue.Description,
		Status:      issue.State.Name,
		Priority:    priorityStr,
		Assignee:    assigneeName,
		CreatedAt:   issue.CreatedAt,
	}
}

// ConvertIssuesToUIModels converts a slice of Linear Issues to domain Issues
func (a *LinearAdapter) ConvertIssuesToUIModels(issues []linear.Issue) []domain.Issue {
	uiIssues := make([]domain.Issue, len(issues))
	for i, issue := range issues {
		uiIssues[i] = a.ConvertIssueToUIModel(issue)
	}
	return uiIssues
}

// ConvertProjectToUIModel converts a Linear Project to a domain Project for UI usage
func (a *LinearAdapter) ConvertProjectToUIModel(project linear.Project) domain.Project {
	return domain.Project{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Status:      project.State,
		Progress:    project.Progress,
		CreatedAt:   parseProjectDate(project.StartDate), // Use start date as created date
	}
}

// ConvertProjectsToUIModels converts a slice of Linear Projects to domain Projects
func (a *LinearAdapter) ConvertProjectsToUIModels(projects []linear.Project) []domain.Project {
	uiProjects := make([]domain.Project, len(projects))
	for i, project := range projects {
		uiProjects[i] = a.ConvertProjectToUIModel(project)
	}
	return uiProjects
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
