package components

import "github.com/linear-tui/linear-tui/internal/ui/mock"

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
	Tickets  []mock.MockTicket
	Projects []mock.MockProject
}

// DataLoadErrorMsg represents a message indicating data loading failed
type DataLoadErrorMsg struct {
	Error error
}
