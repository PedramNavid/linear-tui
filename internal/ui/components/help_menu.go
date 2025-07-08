package components

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyBinding represents a key binding with its description
type KeyBinding struct {
	Key         string
	Description string
}

// KeyBindingSection represents a section of key bindings
type KeyBindingSection struct {
	Title    string
	Bindings []KeyBinding
}

// HelpMenu represents the help menu component
type HelpMenu struct {
	IsVisible bool
	Width     int
	Height    int
	sections  []KeyBindingSection
}

// NewHelpMenu creates a new help menu component
func NewHelpMenu() *HelpMenu {
	return &HelpMenu{
		IsVisible: false,
		sections:  buildKeyBindingSections(),
	}
}

// buildKeyBindingSections creates the key binding sections
func buildKeyBindingSections() []KeyBindingSection {
	return []KeyBindingSection{
		{
			Title: "Global",
			Bindings: []KeyBinding{
				{"q, Ctrl+C", "Quit application"},
				{"?", "Toggle this help menu"},
				{"c", "Create new issue"},
				{"e", "Edit selected issue"},
				{"r", "Refresh data from Linear"},
				{"Ctrl+D", "Toggle detail pane visibility"},
				{"Esc", "Return focus to main pane"},
			},
		},
		{
			Title: "Navigation",
			Bindings: []KeyBinding{
				{"Tab", "Move focus to next pane"},
				{"Shift+Tab", "Move focus to previous pane"},
				{"←→, h/l", "Navigate menu bar items"},
				{"↑↓, k/j", "Navigate list items"},
				{"Enter", "Select/activate item"},
			},
		},
		{
			Title: "Menu Bar",
			Bindings: []KeyBinding{
				{"←→, h/l", "Navigate between Issues and Projects"},
				{"Enter", "Switch to selected view"},
			},
		},
		{
			Title: "Main Pane",
			Bindings: []KeyBinding{
				{"↑↓, k/j", "Navigate through issues/projects"},
				{"Enter", "Select item (updates detail pane)"},
			},
		},
		{
			Title: "Detail Pane",
			Bindings: []KeyBinding{
				{"↑↓, k/j", "Scroll content up/down"},
				{"PgUp/PgDn", "Scroll page up/down"},
				{"Space", "Scroll page down"},
			},
		},
		{
			Title: "Create/Edit Modal",
			Bindings: []KeyBinding{
				{"Tab", "Move to next field"},
				{"Shift+Tab", "Move to previous field"},
				{"↑↓", "Change dropdown selection"},
				{"Enter", "Submit when on Submit button"},
				{"Ctrl+Enter", "Submit from any field"},
				{"Esc", "Cancel and close modal"},
			},
		},
	}
}

// Show displays the help menu
func (h *HelpMenu) Show() {
	h.IsVisible = true
}

// Hide closes the help menu
func (h *HelpMenu) Hide() {
	h.IsVisible = false
}

// Toggle toggles the help menu visibility
func (h *HelpMenu) Toggle() {
	h.IsVisible = !h.IsVisible
}

// Update handles input for the help menu
func (h *HelpMenu) Update(msg tea.Msg) (*HelpMenu, tea.Cmd) {
	if !h.IsVisible {
		return h, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?", "esc":
			h.Hide()
			return h, nil
		case "ctrl+c", "q":
			return h, tea.Quit
		}
	}

	return h, nil
}

// View renders the help menu
func (h *HelpMenu) View(styles *Styles) string {
	if !h.IsVisible {
		return ""
	}

	// Calculate dimensions (80% of terminal size)
	menuWidth := (h.Width * 8) / 10
	menuHeight := (h.Height * 8) / 10

	// Minimum size constraints
	if menuWidth < 70 {
		menuWidth = 70
	}
	if menuHeight < 25 {
		menuHeight = 25
	}

	// Content width (leave space for borders and padding)
	contentWidth := menuWidth - 4

	var content strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(contentWidth).
		Align(lipgloss.Center).
		Bold(true).
		Render("Key Bindings Help")

	content.WriteString(header)
	content.WriteString("\n\n")

	// Render sections in two columns with proper layout
	leftColumn := strings.Builder{}
	rightColumn := strings.Builder{}

	// Split sections into two columns
	sectionsPerColumn := (len(h.sections) + 1) / 2
	leftSections := h.sections[:sectionsPerColumn]
	rightSections := h.sections[sectionsPerColumn:]

	// Calculate column width (leave space for separation)
	columnWidth := (contentWidth - 6) / 2 // 6 chars for spacing between columns
	
	// Render left column
	for i, section := range leftSections {
		if i > 0 {
			leftColumn.WriteString("\n")
		}
		leftColumn.WriteString(h.renderSection(section, columnWidth, styles))
	}

	// Render right column
	for i, section := range rightSections {
		if i > 0 {
			rightColumn.WriteString("\n")
		}
		rightColumn.WriteString(h.renderSection(section, columnWidth, styles))
	}

	// Combine columns using lipgloss JoinHorizontal
	leftContent := leftColumn.String()
	rightContent := rightColumn.String()

	// Use lipgloss to join horizontally with proper spacing
	leftStyle := lipgloss.NewStyle().Width(columnWidth).Align(lipgloss.Left)
	rightStyle := lipgloss.NewStyle().Width(columnWidth).Align(lipgloss.Left)

	combinedContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftContent),
		"      ", // 6 spaces between columns
		rightStyle.Render(rightContent),
	)

	content.WriteString(combinedContent)

	// Footer
	content.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Width(contentWidth).
		Align(lipgloss.Center).
		Render("Press ? or Esc to close")

	content.WriteString(footer)

	// Apply modal styling
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 1).
		Width(menuWidth).
		Height(menuHeight)

	return modalStyle.Render(content.String())
}

// renderSection renders a single key binding section
func (h *HelpMenu) renderSection(section KeyBindingSection, width int, styles *Styles) string {
	var content strings.Builder

	// Section title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(width).
		Bold(true)

	content.WriteString(titleStyle.Render(section.Title))
	content.WriteString("\n")

	// Key bindings with proper formatting
	for _, binding := range section.Bindings {
		keyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD")).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

		// Format with proper spacing: "  Key: Description"
		line := "  " + keyStyle.Render(binding.Key) + ": " + descStyle.Render(binding.Description)

		content.WriteString(line)
		content.WriteString("\n")
	}

	return content.String()
}

// SetDimensions sets the help menu dimensions
func (h *HelpMenu) SetDimensions(width, height int) {
	h.Width = width
	h.Height = height
}

// Helper function for max (already exists in other files, but needed here)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}