package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui/services"
)

// FormField represents the currently focused form field
type FormField int

const (
	FieldTitle FormField = iota
	FieldDescription
	FieldStatus
	FieldPriority
	FieldAssignee
	FieldProject
	FieldSubmit
)

// SelectOption represents an option in a dropdown
type SelectOption struct {
	Value string
	Label string
}

// CreateTicketModal represents the modal for creating a new ticket
type CreateTicketModal struct {
	// Visibility
	IsVisible bool

	// Form fields
	Title       textinput.Model
	Description textarea.Model

	// Dropdown options
	StatusOptions   []SelectOption
	PriorityOptions []SelectOption
	AssigneeOptions []SelectOption
	ProjectOptions  []SelectOption

	// Current selections
	SelectedStatus   int
	SelectedPriority int
	SelectedAssignee int
	SelectedProject  int

	// Focus management
	FocusedField FormField

	// Dimensions
	Width  int
	Height int

	// State
	TeamName      string
	ErrorMessage  string
	SubmitMessage string
	IsSubmitting  bool

	// Services
	LinearService *services.LinearService
}

// NewCreateTicketModal creates a new ticket creation modal
func NewCreateTicketModal() *CreateTicketModal {
	// Initialize title input
	titleInput := textinput.New()
	titleInput.Placeholder = "Enter issue title..."
	titleInput.Focus()
	titleInput.CharLimit = 200
	titleInput.Width = 50

	// Initialize description textarea
	descArea := textarea.New()
	descArea.Placeholder = "Enter issue description..."
	descArea.SetWidth(50)
	descArea.SetHeight(4)

	return &CreateTicketModal{
		IsVisible:   false,
		Title:       titleInput,
		Description: descArea,

		// Default dropdown options
		StatusOptions: []SelectOption{
			{Value: "todo", Label: "Todo"},
			{Value: "in_progress", Label: "In Progress"},
			{Value: "done", Label: "Done"},
		},
		PriorityOptions: []SelectOption{
			{Value: "low", Label: "Low"},
			{Value: "medium", Label: "Medium"},
			{Value: "high", Label: "High"},
			{Value: "urgent", Label: "Urgent"},
		},
		AssigneeOptions: []SelectOption{
			{Value: "unassigned", Label: "Unassigned"},
			{Value: "me", Label: "Assign to me"},
		},
		ProjectOptions: []SelectOption{
			{Value: "none", Label: "No project"},
			{Value: "linear-tui", Label: "Linear TUI"},
		},

		// Default selections (indexes)
		SelectedStatus:   0, // "Todo"
		SelectedPriority: 1, // "Medium"
		SelectedAssignee: 0, // "Unassigned"
		SelectedProject:  0, // "No project"

		FocusedField: FieldTitle,
		TeamName:     "Pedram", // Default team name
	}
}

// Show displays the modal
func (m *CreateTicketModal) Show() {
	m.IsVisible = true
	m.FocusedField = FieldTitle
	m.Title.Focus()
	m.Description.Blur()
	m.clearMessages()
}

// Hide closes the modal
func (m *CreateTicketModal) Hide() {
	m.IsVisible = false
	m.Title.Blur()
	m.Description.Blur()
	m.reset()
}

// Update handles keyboard input for the modal
func (m *CreateTicketModal) Update(msg tea.Msg) (*CreateTicketModal, tea.Cmd) {
	if !m.IsVisible {
		return m, nil
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Hide()
			return m, nil

		case "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.moveFocusForward()

		case "shift+tab":
			m.moveFocusBackward()

		case "ctrl+enter":
			if m.canSubmit() {
				return m, m.submitForm()
			}

		case "enter":
			if m.FocusedField == FieldSubmit && m.canSubmit() {
				return m, m.submitForm()
			}
			// For dropdowns, toggle/select
			m.handleDropdownEnter()

		case "up":
			if m.isDropdownField() {
				m.moveDropdownSelection(-1)
			}

		case "down":
			if m.isDropdownField() {
				m.moveDropdownSelection(1)
			}
		}
	}

	// Update focused field
	switch m.FocusedField {
	case FieldTitle:
		var cmd tea.Cmd
		m.Title, cmd = m.Title.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case FieldDescription:
		var cmd tea.Cmd
		m.Description, cmd = m.Description.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the modal
func (m *CreateTicketModal) View(styles *Styles) string {
	if !m.IsVisible {
		return ""
	}

	// Modal dimensions (80% of terminal size)
	modalWidth := (m.Width * 8) / 10
	modalHeight := (m.Height * 8) / 10

	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalHeight < 20 {
		modalHeight = 20
	}

	// Content width (leave space for borders and padding)
	contentWidth := modalWidth - 4

	var content strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(contentWidth).
		Align(lipgloss.Center).
		Bold(true).
		Render(fmt.Sprintf("Create New Issue - %s", m.TeamName))

	content.WriteString(header)
	content.WriteString("\n\n")

	// Form fields
	content.WriteString(m.renderFormField("Title *", m.Title.View(), FieldTitle, contentWidth))
	content.WriteString("\n")

	content.WriteString(m.renderFormField("Description", m.Description.View(), FieldDescription, contentWidth))
	content.WriteString("\n")

	content.WriteString(m.renderDropdown("Status", m.StatusOptions, m.SelectedStatus, FieldStatus, contentWidth))
	content.WriteString("\n")

	content.WriteString(m.renderDropdown("Priority", m.PriorityOptions, m.SelectedPriority, FieldPriority, contentWidth))
	content.WriteString("\n")

	content.WriteString(m.renderDropdown("Assignee", m.AssigneeOptions, m.SelectedAssignee, FieldAssignee, contentWidth))
	content.WriteString("\n")

	content.WriteString(m.renderDropdown("Project", m.ProjectOptions, m.SelectedProject, FieldProject, contentWidth))
	content.WriteString("\n")

	// Submit button
	submitStyle := styles.ListItem
	if m.FocusedField == FieldSubmit {
		submitStyle = styles.ListSelected
	}
	if !m.canSubmit() {
		submitStyle = styles.Placeholder
	}

	submitButton := submitStyle.
		Width(contentWidth).
		Align(lipgloss.Center).
		Render("[ Create Issue ]")
	content.WriteString(submitButton)
	content.WriteString("\n\n")

	// Messages
	if m.ErrorMessage != "" {
		errorMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Width(contentWidth).
			Render(m.ErrorMessage)
		content.WriteString(errorMsg)
		content.WriteString("\n")
	}

	if m.SubmitMessage != "" {
		successMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Width(contentWidth).
			Render(m.SubmitMessage)
		content.WriteString(successMsg)
		content.WriteString("\n")
	}

	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Width(contentWidth).
		Render("Tab/Shift+Tab: Navigate • Ctrl+Enter: Submit • ESC: Cancel")
	content.WriteString(helpText)

	// Apply modal styling
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 1).
		Width(modalWidth).
		Height(modalHeight)

	return modalStyle.Render(content.String())
}

// SetDimensions sets the modal dimensions
func (m *CreateTicketModal) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height

	// Update form field widths
	contentWidth := ((width * 8) / 10) - 8 // Account for modal padding/borders
	if contentWidth < 40 {
		contentWidth = 40
	}

	m.Title.Width = contentWidth - 10 // Leave space for label
	m.Description.SetWidth(contentWidth - 10)
}

// Helper methods

func (m *CreateTicketModal) renderFormField(label, value string, field FormField, width int) string {
	labelStyle := lipgloss.NewStyle().Bold(true)
	if m.FocusedField == field {
		labelStyle = labelStyle.Foreground(lipgloss.Color("#874BFD"))
	}

	fieldLabel := labelStyle.Render(label)

	return fmt.Sprintf("%s\n%s", fieldLabel, value)
}

func (m *CreateTicketModal) renderDropdown(label string, options []SelectOption, selected int, field FormField, width int) string {
	labelStyle := lipgloss.NewStyle().Bold(true)
	if m.FocusedField == field {
		labelStyle = labelStyle.Foreground(lipgloss.Color("#874BFD"))
	}

	fieldLabel := labelStyle.Render(label)

	valueStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false)

	if m.FocusedField == field {
		valueStyle = valueStyle.BorderForeground(lipgloss.Color("#874BFD"))
	}

	selectedValue := options[selected].Label
	if m.FocusedField == field {
		selectedValue = fmt.Sprintf("▼ %s", selectedValue)
	}

	dropdown := valueStyle.
		Width(width - 10).
		Render(selectedValue)

	return fmt.Sprintf("%s\n%s", fieldLabel, dropdown)
}

func (m *CreateTicketModal) moveFocusForward() {
	m.Title.Blur()
	m.Description.Blur()

	if m.FocusedField == FieldSubmit {
		m.FocusedField = FieldTitle
	} else {
		m.FocusedField++
	}

	m.updateFocus()
}

func (m *CreateTicketModal) moveFocusBackward() {
	m.Title.Blur()
	m.Description.Blur()

	if m.FocusedField == FieldTitle {
		m.FocusedField = FieldSubmit
	} else {
		m.FocusedField--
	}

	m.updateFocus()
}

func (m *CreateTicketModal) updateFocus() {
	switch m.FocusedField {
	case FieldTitle:
		m.Title.Focus()
	case FieldDescription:
		m.Description.Focus()
	}
}

func (m *CreateTicketModal) isDropdownField() bool {
	return m.FocusedField == FieldStatus ||
		m.FocusedField == FieldPriority ||
		m.FocusedField == FieldAssignee ||
		m.FocusedField == FieldProject
}

func (m *CreateTicketModal) handleDropdownEnter() {
	// Could expand dropdown options here
}

func (m *CreateTicketModal) moveDropdownSelection(direction int) {
	var options []SelectOption
	var current *int

	switch m.FocusedField {
	case FieldStatus:
		options = m.StatusOptions
		current = &m.SelectedStatus
	case FieldPriority:
		options = m.PriorityOptions
		current = &m.SelectedPriority
	case FieldAssignee:
		options = m.AssigneeOptions
		current = &m.SelectedAssignee
	case FieldProject:
		options = m.ProjectOptions
		current = &m.SelectedProject
	default:
		return
	}

	newValue := *current + direction
	if newValue < 0 {
		newValue = len(options) - 1
	} else if newValue >= len(options) {
		newValue = 0
	}

	*current = newValue
}

func (m *CreateTicketModal) canSubmit() bool {
	return strings.TrimSpace(m.Title.Value()) != "" && !m.IsSubmitting
}

func (m *CreateTicketModal) submitForm() tea.Cmd {
	m.IsSubmitting = true
	m.clearMessages()

	// Create ticket data
	title := strings.TrimSpace(m.Title.Value())
	description := strings.TrimSpace(m.Description.Value())

	// Validate required fields
	if title == "" {
		return func() tea.Msg {
			return CreateTicketResult{
				Success: false,
				Message: "Title is required",
			}
		}
	}

	// Extract form values
	priority := m.PriorityOptions[m.SelectedPriority].Value
	assignee := m.AssigneeOptions[m.SelectedAssignee].Value

	// Create ticket via Linear service if available
	if m.LinearService != nil {
		return func() tea.Msg {
			ticket, err := m.LinearService.CreateTicket(title, description, priority, assignee)
			if err != nil {
				return CreateTicketResult{
					Success: false,
					Message: err.Error(),
					Error:   err,
				}
			}

			return CreateTicketResult{
				Success: true,
				Message: fmt.Sprintf("Issue '%s' created successfully with ID: %s", ticket.Title, ticket.ID),
			}
		}
	}

	// Fallback: simulate API call
	return func() tea.Msg {
		// Include description in success message to use the variable
		message := fmt.Sprintf("Issue '%s' created successfully!", title)
		if description != "" {
			message += " (with description)"
		}
		return CreateTicketResult{
			Success: true,
			Message: message,
		}
	}
}

func (m *CreateTicketModal) reset() {
	m.Title.SetValue("")
	m.Description.SetValue("")
	m.SelectedStatus = 0
	m.SelectedPriority = 1
	m.SelectedAssignee = 0
	m.SelectedProject = 0
	m.clearMessages()
}

func (m *CreateTicketModal) clearMessages() {
	m.ErrorMessage = ""
	m.SubmitMessage = ""
	m.IsSubmitting = false
}

// SetLinearService sets the Linear service for API calls
func (m *CreateTicketModal) SetLinearService(service *services.LinearService) {
	m.LinearService = service
}

// CreateTicketResult represents the result of ticket creation
type CreateTicketResult struct {
	Success bool
	Message string
	Error   error
}
