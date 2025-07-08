package components

import "github.com/linear-tui/linear-tui/internal/domain"

// AppState represents the current state of the application
type AppState int

const (
	StateLoading AppState = iota
	StateReady
	StateError
	StateRetrying
)

// String returns a string representation of the app state
func (s AppState) String() string {
	switch s {
	case StateLoading:
		return "Loading"
	case StateReady:
		return "Ready"
	case StateError:
		return "Error"
	case StateRetrying:
		return "Retrying"
	default:
		return "Unknown"
	}
}

// LoadingMsg represents a message indicating data is being loaded
type LoadingMsg struct {
	Message string
}

// DataLoadedMsg represents a message indicating data has been loaded successfully
type DataLoadedMsg struct {
	Issues   []domain.Issue
	Projects []domain.Project
}

// DataLoadErrorMsg represents a message indicating data loading failed
type DataLoadErrorMsg struct {
	Error error
}

// RefreshSingleIssueMsg triggers a refresh of a single issue
type RefreshSingleIssueMsg struct {
	IssueID string
}

// SingleIssueRefreshedMsg contains the updated issue data
type SingleIssueRefreshedMsg struct {
	Issue *domain.Issue
	Error error
}
