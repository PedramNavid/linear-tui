package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbletea"
)


// SidebarItem represents a navigation item in the sidebar
type SidebarItem struct {
	Title       string
	Description string
	Key         string
}

// Sidebar represents the sidebar component
type Sidebar struct {
	Items        []SidebarItem
	SelectedItem int
	Width        int
	Height       int
	Focused      bool
}

// NewSidebar creates a new sidebar component
func NewSidebar() *Sidebar {
	return &Sidebar{
		Items: []SidebarItem{
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

// Update handles keyboard input for the sidebar
func (s *Sidebar) Update(msg tea.Msg) (*Sidebar, tea.Cmd) {
	if !s.Focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.SelectedItem > 0 {
				s.SelectedItem--
			}
		case "down", "j":
			if s.SelectedItem < len(s.Items)-1 {
				s.SelectedItem++
			}
		}
	}

	return s, nil
}

// View renders the sidebar
func (s *Sidebar) View(styles *Styles) string {
	var content strings.Builder

	// Title
	title := styles.SidebarTitle.Render("Navigation")
	content.WriteString(title)
	content.WriteString("\n")

	// Items
	for i, item := range s.Items {
		var style lipgloss.Style
		if i == s.SelectedItem && s.Focused {
			style = styles.SidebarSelected
		} else {
			style = styles.SidebarItem
		}
		
		itemText := style.Render(item.Title)
		content.WriteString(itemText)
		content.WriteString("\n")
	}

	// Add some spacing and help text
	content.WriteString("\n")
	helpText := styles.Placeholder.Render("↑/↓ to navigate")
	content.WriteString(helpText)

	// Apply border and sizing
	borderStyle := styles.GetBorderStyle(PaneSidebar, s.getFocusedPane())
	
	// Ensure minimum dimensions (account for border + padding)
	width := s.Width - 2 // 1 char padding on each side
	if width < 1 {
		width = 1
	}
	
	sidebarContent := borderStyle.
		Width(width).
		Render(content.String())

	return sidebarContent
}

// GetSelectedKey returns the key of the currently selected item
func (s *Sidebar) GetSelectedKey() string {
	if s.SelectedItem >= 0 && s.SelectedItem < len(s.Items) {
		return s.Items[s.SelectedItem].Key
	}
	return ""
}

// SetFocus sets the focus state of the sidebar
func (s *Sidebar) SetFocus(focused bool) {
	s.Focused = focused
}

// SetDimensions sets the width and height of the sidebar
func (s *Sidebar) SetDimensions(width, height int) {
	s.Width = width
	s.Height = height
}

// getFocusedPane returns the current focused pane for styling
func (s *Sidebar) getFocusedPane() Pane {
	if s.Focused {
		return PaneSidebar
	}
	return PaneMain // Default to main pane if not focused
}