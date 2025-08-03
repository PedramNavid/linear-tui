package bubble

import (
	"time"
)

// Issue represents a list item for the bubble tea list
type Issue struct {
	ID          string
	title       string
	description string
	status      string
	priority    string
	assignee    string
	createdAt   time.Time
}

func NewItem(id, title, description, status, priority, assignee string, createdAt time.Time) Issue {
	return Issue{
		ID:          id,
		title:       title,
		description: description,
		status:      status,
		priority:    priority,
		assignee:    assignee,
		createdAt:   createdAt,
	}
}

func (i Issue) FilterValue() string {
	return i.title
}

func (i Issue) Title() string {
	return i.title
}

func (i Issue) Description() string {
	return i.description
}
