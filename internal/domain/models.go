package domain

import "time"

// Issue represents a Linear issue in the UI layer
type Issue struct {
	ID          string // Display ID (e.g., "PED-35")
	LinearID    string // Internal Linear ID for API operations
	Title       string
	Description string
	Status      string
	Priority    string
	Assignee    string
	CreatedAt   time.Time
}

// Project represents a Linear project in the UI layer
type Project struct {
	ID          string
	Name        string
	Description string
	Status      string
	Progress    float64
	CreatedAt   time.Time
}
