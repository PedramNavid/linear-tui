package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui/mock"
)

// MainPane represents the main content area
type MainPane struct {
	ViewType     string // "issues" or "projects"
	SelectedItem int
	Width        int
	Height       int
	Focused      bool

	// Viewport state
	ViewportStart int // First visible item index

	// Data
	Tickets  []mock.MockTicket
	Projects []mock.MockProject
}

// NewMainPane creates a new main pane component
func NewMainPane() *MainPane {
	return &MainPane{
		ViewType:     "issues",
		SelectedItem: 0,
		Focused:      false,                // Will be set by layout
		Tickets:      []mock.MockTicket{},  // Start with empty data - will be loaded from Linear
		Projects:     []mock.MockProject{}, // Start with empty data - will be loaded from Linear
	}
}

// Update handles keyboard input for the main pane
func (m *MainPane) Update(msg tea.Msg) (*MainPane, tea.Cmd) {
	if !m.Focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.SelectedItem > 0 {
				m.SelectedItem--
			}
		case "down", "j":
			maxItems := m.getMaxItems()
			if m.SelectedItem < maxItems-1 {
				m.SelectedItem++
			}
		}
	}

	return m, nil
}

// View renders the main pane
func (m *MainPane) View(styles *Styles) string {
	var content strings.Builder

	// Title
	title := m.getTitle()
	titleText := styles.MainTitle.Render(title)
	content.WriteString(titleText)
	content.WriteString("\n")

	// Content based on view type
	if m.ViewType == "issues" {
		m.renderIssuesList(&content, styles)
	} else {
		m.renderProjectsList(&content, styles)
	}

	// Add help text
	content.WriteString("\n")
	helpText := styles.Placeholder.Render("↑/↓ to navigate • Enter to select")
	content.WriteString(helpText)

	// Apply border and sizing
	borderStyle := styles.GetBorderStyle(PaneMain, m.getFocusedPane())

	// Ensure minimum dimensions (account for border + padding)
	width := m.Width - 2 // 1 char padding on each side
	if width < 1 {
		width = 1
	}

	mainContent := borderStyle.
		Width(width).
		Render(content.String())

	return mainContent
}

// renderIssuesList renders the list of issues with viewport logic
func (m *MainPane) renderIssuesList(content *strings.Builder, styles *Styles) {
	// Calculate available height for items (subtract title, help text, and borders)
	availableHeight := m.Height - 6 // Account for title, help, borders (no vertical padding now)
	if availableHeight < 1 {
		availableHeight = 1
	}

	// Calculate available width for content (account for borders and padding)
	availableWidth := m.Width - 6 // Account for borders (2) and padding (4)
	if availableWidth < 10 {
		availableWidth = 10 // Minimum width
	}

	// Calculate visible range around selected item
	startIdx, endIdx := m.calculateVisibleRange(len(m.Tickets), availableHeight)

	for i := startIdx; i <= endIdx && i < len(m.Tickets); i++ {
		ticket := m.Tickets[i]
		var style lipgloss.Style
		if i == m.SelectedItem && m.Focused {
			style = styles.ListSelected
		} else {
			style = styles.ListItem
		}

		// Format: [ID] Title - Status (Priority)
		// Calculate lengths for truncation (without style codes)
		statusText := fmt.Sprintf("(%s)", ticket.Priority)
		metaText := fmt.Sprintf(" - %s %s", ticket.Status, statusText)
		prefixText := fmt.Sprintf("[%s] ", ticket.ID)

		// Calculate space available for title
		reservedSpace := len(prefixText) + len(metaText)
		titleSpace := availableWidth - reservedSpace

		// Truncate title if necessary
		title := ticket.Title
		var finalMetaText string

		if len(title) > titleSpace && titleSpace > 3 {
			title = title[:titleSpace-3] + "..."
			finalMetaText = metaText
		} else if titleSpace <= 3 {
			// Very limited space, just show prefix and truncated title
			if availableWidth > len(prefixText)+3 {
				title = title[:availableWidth-len(prefixText)-3] + "..."
				finalMetaText = "" // Remove metadata if no space
			} else {
				title = "..."
				finalMetaText = ""
			}
		} else {
			finalMetaText = metaText
		}

		// Build the final text with proper styling for priority
		var itemText string
		if finalMetaText != "" {
			// Apply priority styling to just the priority part
			priorityStyled := styles.GetStatusStyle(ticket.Priority).Render(ticket.Priority)
			itemText = fmt.Sprintf("%s%s - %s (%s)", prefixText, title, ticket.Status, priorityStyled)
		} else {
			itemText = prefixText + title
		}

		renderedItem := style.Render(itemText)
		content.WriteString(renderedItem)
		content.WriteString("\n")
	}
}

// renderProjectsList renders the list of projects with viewport logic
func (m *MainPane) renderProjectsList(content *strings.Builder, styles *Styles) {
	// Calculate available height for items (subtract title, help text, and borders)
	availableHeight := m.Height - 6 // Account for title, help, borders (no vertical padding now)
	if availableHeight < 1 {
		availableHeight = 1
	}

	// Calculate visible range around selected item
	startIdx, endIdx := m.calculateVisibleRange(len(m.Projects), availableHeight)

	for i := startIdx; i <= endIdx && i < len(m.Projects); i++ {
		project := m.Projects[i]
		var style lipgloss.Style
		if i == m.SelectedItem && m.Focused {
			style = styles.ListSelected
		} else {
			style = styles.ListItem
		}

		// Format: [ID] Name - Status (Progress%)
		progressText := fmt.Sprintf("%.0f%%", project.Progress*100)
		itemText := fmt.Sprintf("[%s] %s - %s (%s)",
			project.ID,
			project.Name,
			project.Status,
			progressText)

		renderedItem := style.Render(itemText)
		content.WriteString(renderedItem)
		content.WriteString("\n")
	}
}

// getTitle returns the title for the current view
func (m *MainPane) getTitle() string {
	switch m.ViewType {
	case "issues":
		return fmt.Sprintf("Issues (%d)", len(m.Tickets))
	case "projects":
		return fmt.Sprintf("Projects (%d)", len(m.Projects))
	default:
		return "Main"
	}
}

// getMaxItems returns the maximum number of items in the current view
func (m *MainPane) getMaxItems() int {
	switch m.ViewType {
	case "issues":
		return len(m.Tickets)
	case "projects":
		return len(m.Projects)
	default:
		return 0
	}
}

// SetViewType changes the view type and resets selection
func (m *MainPane) SetViewType(viewType string) {
	m.ViewType = viewType
	m.SelectedItem = 0
	m.ViewportStart = 0 // Reset viewport when switching views
}

// GetSelectedTicket returns the currently selected ticket
func (m *MainPane) GetSelectedTicket() *mock.MockTicket {
	if m.ViewType == "issues" && m.SelectedItem >= 0 && m.SelectedItem < len(m.Tickets) {
		return &m.Tickets[m.SelectedItem]
	}
	return nil
}

// GetSelectedProject returns the currently selected project
func (m *MainPane) GetSelectedProject() *mock.MockProject {
	if m.ViewType == "projects" && m.SelectedItem >= 0 && m.SelectedItem < len(m.Projects) {
		return &m.Projects[m.SelectedItem]
	}
	return nil
}

// SetFocus sets the focus state of the main pane
func (m *MainPane) SetFocus(focused bool) {
	m.Focused = focused
}

// SetDimensions sets the width and height of the main pane
func (m *MainPane) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

// getFocusedPane returns the current focused pane for styling
func (m *MainPane) getFocusedPane() Pane {
	if m.Focused {
		return PaneMain
	}
	return PaneMenu // Default if not focused
}

// calculateVisibleRange calculates which items should be visible in the viewport
func (m *MainPane) calculateVisibleRange(totalItems, availableHeight int) (int, int) {
	if totalItems == 0 || availableHeight <= 0 {
		return 0, 0
	}

	// If all items fit, show them all
	if totalItems <= availableHeight {
		m.ViewportStart = 0
		return 0, totalItems - 1
	}

	// Only adjust viewport if selected item goes out of view
	viewportEnd := m.ViewportStart + availableHeight - 1

	// Scroll down if selected item is below viewport
	if m.SelectedItem > viewportEnd {
		m.ViewportStart = m.SelectedItem - availableHeight + 1
	}

	// Scroll up if selected item is above viewport
	if m.SelectedItem < m.ViewportStart {
		m.ViewportStart = m.SelectedItem
	}

	// Ensure viewport stays within bounds
	if m.ViewportStart < 0 {
		m.ViewportStart = 0
	}
	if m.ViewportStart+availableHeight > totalItems {
		m.ViewportStart = totalItems - availableHeight
		if m.ViewportStart < 0 {
			m.ViewportStart = 0
		}
	}

	endIdx := m.ViewportStart + availableHeight - 1
	if endIdx >= totalItems {
		endIdx = totalItems - 1
	}

	return m.ViewportStart, endIdx
}
