package bubble

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the main application state
type Model struct {
	Title  string
	Styles Styles

	issues        list.Model
	cursor        int
	selectedIssue map[int]struct{}
	err           error
	width         int
	height        int
	style         lipgloss.Style
	onStartup     bool
}

func (m Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.updatePagination()
}

func (m Model) updatePagination() {
}

func (m Model) SetOnStartup(onStartup bool) {
	m.onStartup = onStartup
}

// NewModel creates a new model instance
func NewModel() Model {
	return Model{
		issues:        list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		cursor:        0,
		selectedIssue: make(map[int]struct{}),
		onStartup:     true, // Initialize to true so test data loads on startup
	}
}

// Update handles incoming messages and updates the model accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	windowSizeMsg, isWindowSizeMsg := msg.(tea.WindowSizeMsg)
	if m.onStartup && !isWindowSizeMsg {
		return m, nil
	}

	if m.onStartup && isWindowSizeMsg {
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.SetSize(windowSizeMsg.Width-h, windowSizeMsg.Height-v)
		m.SetOnStartup(false)
		m.issues.SetItems(m.loadTestData())
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.issues.SetSize(windowSizeMsg.Width-h, windowSizeMsg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			m.issues.CursorUp()

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			m.issues.CursorDown()

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selectedIssue[m.cursor]
			if ok {
				delete(m.selectedIssue, m.cursor)
			} else {
				m.selectedIssue[m.cursor] = struct{}{}
			}
		}

	case ErrMsg:
		m.err = msg
		return m, tea.Quit

	case IssuesMsg:
		// Convert domain.Issue to list.Item
		items := make([]list.Item, len(msg.Issues()))
		for i, issue := range msg.Issues() {
			items[i] = Item{
				ID:          issue.ID,
				Title:       issue.Title,
				Description: issue.Description,
				Status:      issue.Status,
				Priority:    issue.Priority,
				Assignee:    issue.Assignee,
				CreatedAt:   issue.CreatedAt,
			}
		}
		m.issues.SetItems(items)
		return m, nil
	}

	var cmd tea.Cmd
	m.issues, cmd = m.issues.Update(msg)
	return m, cmd

}

// Init initializes the model and returns the initial command
func (m Model) Init() tea.Cmd {
	return nil
}

// loadTestData loads test data for development
func (m Model) loadTestData() []list.Item {
	// Create test items directly to avoid import cycles
	testItems := []list.Item{
		Item{
			ID:          "LIN-001",
			Title:       "Implement user authentication",
			Description: "Add OAuth2 authentication flow for user login. This should include Google and GitHub providers with proper session management and token refresh capabilities.",
			Status:      "In Progress",
			Priority:    "High",
			Assignee:    "Alice Johnson",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
		Item{
			ID:          "LIN-002",
			Title:       "Fix database connection pooling",
			Description: "The current database connection pool is not properly handling timeouts and connection limits. This causes the application to hang under high load.",
			Status:      "Todo",
			Priority:    "High",
			Assignee:    "Bob Smith",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		Item{
			ID:          "LIN-003",
			Title:       "Add dark mode support",
			Description: "Implement dark mode theme throughout the application with proper color scheme and user preference storage.",
			Status:      "Done",
			Priority:    "Medium",
			Assignee:    "Carol Davis",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		},
		Item{
			ID:          "LIN-004",
			Title:       "Optimize API response times",
			Description: "Current API response times are averaging 500ms. Need to implement caching strategies and query optimization to bring this down to under 200ms.",
			Status:      "In Progress",
			Priority:    "High",
			Assignee:    "David Wilson",
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
		Item{
			ID:          "LIN-005",
			Title:       "Create user onboarding flow",
			Description: "Design and implement a comprehensive onboarding experience for new users including tutorials and tooltips.",
			Status:      "Todo",
			Priority:    "Medium",
			Assignee:    "Unassigned",
			CreatedAt:   time.Now().Add(-96 * time.Hour),
		},
	}

	return testItems
}
