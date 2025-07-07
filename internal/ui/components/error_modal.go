package components

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrorModal represents a modal for displaying errors with action options
type ErrorModal struct {
	IsVisible      bool
	Title          string
	ErrorMessage   string
	SelectedAction int
	Actions        []ErrorAction
	Width          int
	Height         int
}

// ErrorAction represents an action the user can take in response to an error
type ErrorAction struct {
	Label string
	Key   string
}

// ErrorModalResult represents the result of user action in error modal
type ErrorModalResult struct {
	Action string // "retry", "quit", etc.
}

// NewErrorModal creates a new error modal
func NewErrorModal() *ErrorModal {
	return &ErrorModal{
		IsVisible:      false,
		SelectedAction: 0,
		Actions: []ErrorAction{
			{Label: "Retry", Key: "retry"},
			{Label: "Quit", Key: "quit"},
		},
	}
}

// Show displays the error modal with the given title and message
func (m *ErrorModal) Show(title, message string) {
	m.IsVisible = true
	m.Title = title
	m.ErrorMessage = message
	m.SelectedAction = 0 // Default to first action (Retry)
}

// Hide closes the error modal
func (m *ErrorModal) Hide() {
	m.IsVisible = false
	m.Title = ""
	m.ErrorMessage = ""
	m.SelectedAction = 0
}

// Update handles keyboard input for the error modal
func (m *ErrorModal) Update(msg tea.Msg) (*ErrorModal, tea.Cmd) {
	if !m.IsVisible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.SelectedAction > 0 {
				m.SelectedAction--
			} else {
				m.SelectedAction = len(m.Actions) - 1
			}

		case "right", "l":
			if m.SelectedAction < len(m.Actions)-1 {
				m.SelectedAction++
			} else {
				m.SelectedAction = 0
			}

		case "enter", " ":
			if m.SelectedAction < len(m.Actions) {
				action := m.Actions[m.SelectedAction].Key
				m.Hide()
				return m, func() tea.Msg {
					return ErrorModalResult{Action: action}
				}
			}

		case "r":
			// Quick retry shortcut
			m.Hide()
			return m, func() tea.Msg {
				return ErrorModalResult{Action: "retry"}
			}

		case "q", "ctrl+c":
			// Quick quit shortcut
			m.Hide()
			return m, func() tea.Msg {
				return ErrorModalResult{Action: "quit"}
			}
		}
	}

	return m, nil
}

// View renders the error modal
func (m *ErrorModal) View(styles *Styles) string {
	if !m.IsVisible {
		return ""
	}

	// Modal dimensions (60% of terminal size, centered)
	modalWidth := (m.Width * 6) / 10
	modalHeight := (m.Height * 6) / 10

	if modalWidth < 50 {
		modalWidth = 50
	}
	if modalHeight < 15 {
		modalHeight = 15
	}

	// Content width (leave space for borders and padding)
	contentWidth := modalWidth - 4

	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF4444")).
		Bold(true).
		Padding(0, 1).
		Width(contentWidth).
		Align(lipgloss.Center)

	content.WriteString(titleStyle.Render(m.Title))
	content.WriteString("\n\n")

	// Error message
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Width(contentWidth).
		Align(lipgloss.Left)

	// Word wrap the error message
	wrappedMessage := wordWrap(m.ErrorMessage, contentWidth-4)
	content.WriteString(messageStyle.Render(wrappedMessage))
	content.WriteString("\n\n")

	// Actions (buttons)
	var actionButtons strings.Builder
	for i, action := range m.Actions {
		var buttonStyle lipgloss.Style
		if i == m.SelectedAction {
			buttonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)
		} else {
			buttonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Background(lipgloss.Color("#444444")).
				Padding(0, 2).
				Margin(0, 1)
		}

		actionButtons.WriteString(buttonStyle.Render(action.Label))
	}

	// Center the action buttons
	buttonsStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center)

	content.WriteString(buttonsStyle.Render(actionButtons.String()))
	content.WriteString("\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(contentWidth).
		Align(lipgloss.Center)

	helpText := "Use ←/→ or h/l to navigate • Enter to select • r to retry • q to quit"
	content.WriteString(helpStyle.Render(helpText))

	// Create modal border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF4444")).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight)

	modal := modalStyle.Render(content.String())

	// Create overlay background
	overlayStyle := lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height).
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Align(lipgloss.Center, lipgloss.Center)

	return overlayStyle.Render(modal)
}

// SetDimensions sets the modal dimensions
func (m *ErrorModal) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
}
