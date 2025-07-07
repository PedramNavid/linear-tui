package linear

import "time"

type Issue struct {
	ID          string     `json:"id"`
	Identifier  string     `json:"identifier"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       IssueState `json:"state"`
	Priority    int        `json:"priority"`
	Assignee    *User      `json:"assignee"`
	Team        *Team      `json:"team"`
	Project     *Project   `json:"project"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type IssueState struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

type Project struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	State       string  `json:"state"`
	Progress    float64 `json:"progress"`
	StartDate   *string `json:"startDate"`
	TargetDate  *string `json:"targetDate"`
}

type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

type Comment struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	User      *User     `json:"user"`
	Issue     *Issue    `json:"issue"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateIssueInput struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	TeamID      string `json:"teamId"`
	Priority    int    `json:"priority,omitempty"`
	AssigneeID  string `json:"assigneeId,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	StateID     string `json:"stateId,omitempty"`
}

type UpdateIssueInput struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	AssigneeID  string `json:"assigneeId,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	StateID     string `json:"stateId,omitempty"`
}

type IssuesResponse struct {
	Issues struct {
		Nodes    []Issue  `json:"nodes"`
		PageInfo PageInfo `json:"pageInfo"`
	} `json:"issues"`
}

type ProjectsResponse struct {
	Projects struct {
		Nodes    []Project `json:"nodes"`
		PageInfo PageInfo  `json:"pageInfo"`
	} `json:"projects"`
}

type TeamsResponse struct {
	Teams struct {
		Nodes    []Team   `json:"nodes"`
		PageInfo PageInfo `json:"pageInfo"`
	} `json:"teams"`
}

type UsersResponse struct {
	Users struct {
		Nodes    []User   `json:"nodes"`
		PageInfo PageInfo `json:"pageInfo"`
	} `json:"users"`
}

type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
}
