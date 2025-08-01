package bubble

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// Item represents a list item for the bubble tea list
type Item struct {
	ID          string
	Title       string
	Description string
	Status      string
	Priority    string
	Assignee    string
	CreatedAt   time.Time
}

// FilterValue implements list.Item interface
func (i Item) FilterValue() string {
	return i.Title
}
