package bubble

import (
	"github.com/linear-tui/linear-tui/internal/domain"
)

// IssuesMsg represents a message containing issues data
type IssuesMsg struct {
	issues []domain.Issue
}

// Issues returns the issues from the message
func (m IssuesMsg) Issues() []domain.Issue { return m.issues }

// ErrMsg represents an error message
type ErrMsg struct{ err error }

// Error returns the error string
func (e ErrMsg) Error() string { return e.err.Error() }
