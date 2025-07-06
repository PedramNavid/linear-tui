package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/linear-tui/linear-tui/internal/ui/mock"
)

// DetailPane represents the detail view pane
type DetailPane struct {
	Width   int
	Height  int
	Focused bool
	
	// Content
	SelectedTicket  *mock.MockTicket
	SelectedProject *mock.MockProject
	ViewType        string // "issues" or "projects"
}

// NewDetailPane creates a new detail pane component
func NewDetailPane() *DetailPane {
	return &DetailPane{
		Focused:  false,
		ViewType: "issues",
	}
}

// Update handles keyboard input for the detail pane
func (d *DetailPane) Update(msg tea.Msg) (*DetailPane, tea.Cmd) {
	// Detail pane is mostly read-only for now
	// Future keyboard navigation for scrolling can be added here
	return d, nil
}

// View renders the detail pane
func (d *DetailPane) View(styles *Styles) string {
	var content strings.Builder

	// Title
	title := styles.DetailTitle.Render("Details")
	content.WriteString(title)
	content.WriteString("\n")

	// Content based on what's selected
	if d.ViewType == "issues" && d.SelectedTicket != nil {
		d.renderTicketDetails(&content, styles)
	} else if d.ViewType == "projects" && d.SelectedProject != nil {
		d.renderProjectDetails(&content, styles)
	} else {
		d.renderEmptyState(&content, styles)
	}

	// Add placeholders for future features
	content.WriteString("\n")
	content.WriteString(styles.Placeholder.Render("--- Placeholders for future features ---"))
	content.WriteString("\n")
	content.WriteString(styles.Placeholder.Render("• Comments"))
	content.WriteString("\n")
	content.WriteString(styles.Placeholder.Render("• Sub-issues"))
	content.WriteString("\n")
	content.WriteString(styles.Placeholder.Render("• Attachments"))
	content.WriteString("\n")
	content.WriteString(styles.Placeholder.Render("• Activity log"))

	// Apply border and sizing
	borderStyle := styles.GetBorderStyle(PaneDetail, d.getFocusedPane())
	
	// Ensure minimum dimensions (account for border + padding)
	width := d.Width - 2 // 1 char padding on each side
	if width < 1 {
		width = 1
	}
	
	detailContent := borderStyle.
		Width(width).
		Render(content.String())

	return detailContent
}

// renderTicketDetails renders the details of a selected ticket
func (d *DetailPane) renderTicketDetails(content *strings.Builder, styles *Styles) {
	ticket := d.SelectedTicket
	
	// Title
	titleText := styles.DetailContent.Render(fmt.Sprintf("Title: %s", ticket.Title))
	content.WriteString(titleText)
	content.WriteString("\n\n")
	
	// Metadata
	metaLines := []string{
		fmt.Sprintf("ID: %s", ticket.ID),
		fmt.Sprintf("Status: %s", ticket.Status),
		fmt.Sprintf("Priority: %s", ticket.Priority),
		fmt.Sprintf("Assignee: %s", ticket.Assignee),
		fmt.Sprintf("Created: %s", ticket.CreatedAt.Format("2006-01-02 15:04")),
	}
	
	for _, line := range metaLines {
		metaText := styles.DetailMeta.Render(line)
		content.WriteString(metaText)
		content.WriteString("\n")
	}
	
	content.WriteString("\n")
	
	// Description
	descTitle := styles.DetailContent.Render("Description:")
	content.WriteString(descTitle)
	content.WriteString("\n")
	
	// Word wrap the description
	wrappedDesc := wordWrap(ticket.Description, d.Width-8) // Account for padding
	descContent := styles.DetailContent.Render(wrappedDesc)
	content.WriteString(descContent)
}

// renderProjectDetails renders the details of a selected project
func (d *DetailPane) renderProjectDetails(content *strings.Builder, styles *Styles) {
	project := d.SelectedProject
	
	// Title
	titleText := styles.DetailContent.Render(fmt.Sprintf("Project: %s", project.Name))
	content.WriteString(titleText)
	content.WriteString("\n\n")
	
	// Metadata
	metaLines := []string{
		fmt.Sprintf("ID: %s", project.ID),
		fmt.Sprintf("Status: %s", project.Status),
		fmt.Sprintf("Progress: %.0f%%", project.Progress*100),
		fmt.Sprintf("Created: %s", project.CreatedAt.Format("2006-01-02 15:04")),
	}
	
	for _, line := range metaLines {
		metaText := styles.DetailMeta.Render(line)
		content.WriteString(metaText)
		content.WriteString("\n")
	}
	
	content.WriteString("\n")
	
	// Description
	descTitle := styles.DetailContent.Render("Description:")
	content.WriteString(descTitle)
	content.WriteString("\n")
	
	// Word wrap the description
	wrappedDesc := wordWrap(project.Description, d.Width-8) // Account for padding
	descContent := styles.DetailContent.Render(wrappedDesc)
	content.WriteString(descContent)
}

// renderEmptyState renders the empty state when nothing is selected
func (d *DetailPane) renderEmptyState(content *strings.Builder, styles *Styles) {
	emptyText := styles.Placeholder.Render("Select an item to view details")
	content.WriteString(emptyText)
}

// SetSelectedTicket sets the selected ticket for display
func (d *DetailPane) SetSelectedTicket(ticket *mock.MockTicket) {
	d.SelectedTicket = ticket
	d.SelectedProject = nil
	d.ViewType = "issues"
}

// SetSelectedProject sets the selected project for display
func (d *DetailPane) SetSelectedProject(project *mock.MockProject) {
	d.SelectedProject = project
	d.SelectedTicket = nil
	d.ViewType = "projects"
}

// SetFocus sets the focus state of the detail pane
func (d *DetailPane) SetFocus(focused bool) {
	d.Focused = focused
}

// SetDimensions sets the width and height of the detail pane
func (d *DetailPane) SetDimensions(width, height int) {
	d.Width = width
	d.Height = height
}

// getFocusedPane returns the current focused pane for styling
func (d *DetailPane) getFocusedPane() Pane {
	if d.Focused {
		return PaneDetail
	}
	return PaneMain // Default if not focused
}

// wordWrap wraps text to fit within a given width
func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	var result strings.Builder
	var currentLine strings.Builder
	
	for _, word := range words {
		// Check if adding this word would exceed the width
		if currentLine.Len() > 0 && currentLine.Len()+len(word)+1 > width {
			result.WriteString(currentLine.String())
			result.WriteString("\n")
			currentLine.Reset()
		}
		
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	
	// Add any remaining content
	if currentLine.Len() > 0 {
		result.WriteString(currentLine.String())
	}
	
	return result.String()
}