package bubble

import (
	"time"
)

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

func NewItem(id, title, description, status, priority, assignee string, createdAt time.Time) Item {
	return Item{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		Assignee:    assignee,
		CreatedAt:   createdAt,
	}
}

func (i Item) FilterValue() string {
	return i.Title
}
