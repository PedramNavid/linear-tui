package components

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuBarItem represents a navigation item in the menu bar
type MenuBarItem struct {
	Title       string
	Description string
	Key         string
}

// MenuBar represents the top menu bar component
type MenuBar struct {
	Items        []MenuBarItem
	SelectedItem int
	Width        int
	Height       int
	Focused      bool
}

// NewMenuBar creates a new menu bar component
func NewMenuBar() *MenuBar {
	return &MenuBar{
		Items: []MenuBarItem{
			{
				Title:       "Issues",
				Description: "View and manage issues",
				Key:         "issues",
			},
			{
				Title:       "Projects",
				Description: "View and manage projects",
				Key:         "projects",
			},
		},
		SelectedItem: 0,
		Focused:      false,
	}
}

// Update handles keyboard input for the menu bar
func (m *MenuBar) Update(msg tea.Msg) (*MenuBar, tea.Cmd) {
	if !m.Focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.SelectedItem > 0 {
				m.SelectedItem--
			}
		case "right", "l":
			if m.SelectedItem < len(m.Items)-1 {
				m.SelectedItem++
			}
		}
	}

	return m, nil
}

// View renders the menu bar
func (m *MenuBar) View(styles *Styles) string {
	var content strings.Builder

	// Create horizontal navigation with items separated by " | "
	for i, item := range m.Items {
		var style lipgloss.Style
		if i == m.SelectedItem && m.Focused {
			style = styles.MenuSelected
		} else {
			style = styles.MenuItem
		}

		itemText := style.Render(item.Title)
		content.WriteString(itemText)

		// Add separator if not the last item
		if i < len(m.Items)-1 {
			separator := styles.MenuItem.Render(" | ")
			content.WriteString(separator)
		}
	}

	// Add help text on a new line
	content.WriteString("\n")
	helpText := styles.Placeholder.Render("←/→ to navigate • Enter to select • Tab to switch panes")
	content.WriteString(helpText)

	// Apply border and sizing
	borderStyle := styles.GetBorderStyle(PaneMenu, m.getFocusedPane())

	// Ensure minimum dimensions (account for border + padding)
	width := m.Width - 2 // 1 char padding on each side
	if width < 1 {
		width = 1
	}

	menuContent := borderStyle.
		Width(width).
		Render(content.String())

	return menuContent
}

// GetSelectedKey returns the key of the currently selected item
func (m *MenuBar) GetSelectedKey() string {
	if m.SelectedItem >= 0 && m.SelectedItem < len(m.Items) {
		return m.Items[m.SelectedItem].Key
	}
	return ""
}

// SetFocus sets the focus state of the menu bar
func (m *MenuBar) SetFocus(focused bool) {
	m.Focused = focused
}

// SetDimensions sets the width and height of the menu bar
func (m *MenuBar) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

// getFocusedPane returns the current focused pane for styling
func (m *MenuBar) getFocusedPane() Pane {
	if m.Focused {
		return PaneMenu
	}
	return PaneMain // Default to main pane if not focused
}
