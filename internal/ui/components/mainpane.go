package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/domain"
)

const (
	// Layout constants (matching DetailPane for consistency)
	mainBorderWidth       = 2 // Border takes 1 char on each side
	mainHorizontalPadding = 4 // Padding: 2 chars on each side (from border style)

	// Fixed element heights
	mainTitleHeight = 2 // Title line + margin
	mainHelpHeight  = 1 // Help text line
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
	Issues   []domain.Issue
	Projects []domain.Project
}

// NewMainPane creates a new main pane component
func NewMainPane() *MainPane {
	return &MainPane{
		ViewType:     "issues",
		SelectedItem: 0,
		Focused:      false,              // Will be set by layout
		Issues:       []domain.Issue{},   // Start with empty data - will be loaded from Linear
		Projects:     []domain.Project{}, // Start with empty data - will be loaded from Linear
	}
}

// getContentWidth returns the available width for content after accounting for borders and padding
func (m *MainPane) getContentWidth() int {
	contentWidth := m.Width - mainBorderWidth - mainHorizontalPadding
	if contentWidth < 1 {
		return 1
	}
	return contentWidth
}

// getAvailableHeight returns the available height for list items after accounting for fixed elements
func (m *MainPane) getAvailableHeight() int {
	// Account for: border, title, help text
	totalFixedHeight := mainBorderWidth + mainTitleHeight + mainHelpHeight
	availableHeight := m.Height - totalFixedHeight
	if availableHeight < 1 {
		return 1
	}
	return availableHeight
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
	// Build sections
	titleSection := m.buildTitleSection(styles)
	listSection := m.buildListSection(styles)
	helpSection := m.buildHelpSection(styles)

	// Compose sections using JoinVertical
	content := lipgloss.JoinVertical(
		lipgloss.Top,
		titleSection,
		listSection,
		helpSection,
	)

	// Apply border and sizing
	borderStyle := styles.GetBorderStyle(PaneMain, m.getFocusedPane())
	width := max(m.Width-mainBorderWidth, 1)

	return borderStyle.Width(width).Render(content)
}

// buildTitleSection builds the title section
func (m *MainPane) buildTitleSection(styles *Styles) string {
	title := m.getTitle()
	return styles.MainTitle.Width(m.getContentWidth()).Render(title)
}

// buildListSection builds the list content section
func (m *MainPane) buildListSection(styles *Styles) string {
	var content strings.Builder

	if m.ViewType == "issues" {
		m.renderIssuesListStyled(&content, styles)
	} else {
		m.renderProjectsListStyled(&content, styles)
	}

	// Ensure consistent height by padding if needed
	lines := strings.Split(content.String(), "\n")
	availableHeight := m.getAvailableHeight()

	// Remove empty last line if present
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Pad to fill available height
	for len(lines) < availableHeight {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// buildHelpSection builds the help text section
func (m *MainPane) buildHelpSection(styles *Styles) string {
	return styles.Placeholder.Width(m.getContentWidth()).Render("↑/↓ to navigate • Enter to select • Ctrl+D to toggle details pane • ? for help")
}

// getTitle returns the title for the current view
func (m *MainPane) getTitle() string {
	switch m.ViewType {
	case "issues":
		return fmt.Sprintf("Issues (%d)", len(m.Issues))
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
		return len(m.Issues)
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

// GetSelectedTicket returns the currently selected issue
func (m *MainPane) GetSelectedIssue() *domain.Issue {
	if m.ViewType == "issues" && m.SelectedItem >= 0 && m.SelectedItem < len(m.Issues) {
		return &m.Issues[m.SelectedItem]
	}
	return nil
}

// GetSelectedProject returns the currently selected project
func (m *MainPane) GetSelectedProject() *domain.Project {
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

// renderIssuesListStyled renders the list of issues with pre-applied styles
func (m *MainPane) renderIssuesListStyled(content *strings.Builder, styles *Styles) {
	availableHeight := m.getAvailableHeight()
	contentWidth := m.getContentWidth()

	// Calculate visible range
	startIdx, endIdx := m.calculateVisibleRange(len(m.Issues), availableHeight)

	for i := startIdx; i <= endIdx && i < len(m.Issues); i++ {
		if i > startIdx {
			content.WriteString("\n")
		}

		issue := m.Issues[i]
		isSelected := i == m.SelectedItem && m.Focused

		// Build the item text
		itemText := m.formatIssueItem(issue, contentWidth, styles)

		// Apply selection style
		var style lipgloss.Style
		if isSelected {
			style = styles.ListSelected.Width(contentWidth).MaxWidth(contentWidth)
		} else {
			style = styles.ListItem.Width(contentWidth).MaxWidth(contentWidth)
		}

		content.WriteString(style.Render(itemText))
	}
}

// formatTicketItem formats a issue item with proper styling
func (m *MainPane) formatIssueItem(issue domain.Issue, width int, styles *Styles) string {
	// Build components
	id := fmt.Sprintf("[%s]", issue.ID)
	status := issue.Status
	priority := fmt.Sprintf("(%s)", issue.Priority)

	// Apply priority color
	priorityStyled := styles.GetStatusStyle(issue.Priority).Render(priority)

	// Calculate available space for title
	metadataLen := len(id) + 1 + len(status) + 1 + len(priority) + 3 // spaces and separators
	titleSpace := width - metadataLen

	// Use lipgloss to handle title truncation
	titleStyle := lipgloss.NewStyle().MaxWidth(titleSpace)
	titleRendered := titleStyle.Render(issue.Title)

	// Build final string
	return fmt.Sprintf("%s %s - %s %s", id, titleRendered, status, priorityStyled)
}

// renderProjectsListStyled renders the list of projects with pre-applied styles
func (m *MainPane) renderProjectsListStyled(content *strings.Builder, styles *Styles) {
	availableHeight := m.getAvailableHeight()
	contentWidth := m.getContentWidth()

	// Calculate visible range
	startIdx, endIdx := m.calculateVisibleRange(len(m.Projects), availableHeight)

	for i := startIdx; i <= endIdx && i < len(m.Projects); i++ {
		if i > startIdx {
			content.WriteString("\n")
		}

		project := m.Projects[i]
		isSelected := i == m.SelectedItem && m.Focused

		// Build the item text
		itemText := m.formatProjectItem(project, contentWidth)

		// Apply selection style
		var style lipgloss.Style
		if isSelected {
			style = styles.ListSelected.Width(contentWidth).MaxWidth(contentWidth)
		} else {
			style = styles.ListItem.Width(contentWidth).MaxWidth(contentWidth)
		}

		content.WriteString(style.Render(itemText))
	}
}

// formatProjectItem formats a project item
func (m *MainPane) formatProjectItem(project domain.Project, width int) string {
	// Build components
	id := fmt.Sprintf("[%s]", project.ID)
	status := project.Status
	progress := fmt.Sprintf("(%.0f%%)", project.Progress*100)

	// Calculate available space for name
	metadataLen := len(id) + 1 + len(status) + 1 + len(progress) + 3 // spaces and separators
	nameSpace := width - metadataLen

	// Use lipgloss to handle name truncation
	nameStyle := lipgloss.NewStyle().MaxWidth(nameSpace)
	nameRendered := nameStyle.Render(project.Name)

	// Build final string
	return fmt.Sprintf("%s %s - %s %s", id, nameRendered, status, progress)
}

// UpdateSingleIssue updates a single issue in the list
func (m *MainPane) UpdateSingleIssue(updatedIssue domain.Issue) {
	for i := range m.Issues {
		if m.Issues[i].LinearID == updatedIssue.LinearID {
			// Update the issue while preserving the position
			m.Issues[i] = updatedIssue
			break
		}
	}
}
