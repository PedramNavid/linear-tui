package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
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

	// Viewport for scrolling
	viewport viewport.Model
	ready    bool // Whether viewport has been initialized
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
	var cmd tea.Cmd
	
	switch msg.(type) {
	case tea.WindowSizeMsg:
		// Just update dimensions if viewport is already ready
		if d.ready {
			d.updateViewportDimensions()
			// Re-set content to reflow for new dimensions
			d.viewport.SetContent(d.buildContent())
		}
	}
	
	// Update viewport if ready
	if d.ready {
		// Only process viewport updates when focused for keyboard events
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if d.Focused {
				d.viewport, cmd = d.viewport.Update(msg)
			}
		} else {
			// Non-keyboard messages always go to viewport
			d.viewport, cmd = d.viewport.Update(msg)
		}
	}
	
	return d, cmd
}

// View renders the detail pane
func (d *DetailPane) View(styles *Styles) string {
	var content strings.Builder

	// Title with scroll indicator and focus hint
	title := "Details"
	if d.ready && d.viewport.TotalLineCount() > d.viewport.Height {
		scrollPercent := d.viewport.ScrollPercent()
		if scrollPercent == 0 {
			title += " (Top)"
		} else if scrollPercent >= 0.99 {
			title += " (Bottom)"
		} else {
			title += fmt.Sprintf(" (%d%%)", int(scrollPercent*100))
		}
	}
	
	// Add focus indicator and keyboard hints
	if d.Focused {
		title += " [FOCUSED - j/k or ↑/↓: scroll]"
	}
	
	
	content.WriteString(styles.DetailTitle.Render(title))
	content.WriteString("\n")

	// Viewport content or loading state
	if !d.ready {
		content.WriteString(styles.Placeholder.Render("Loading..."))
	} else {
		// Get viewport content
		viewportContent := d.viewport.View()
		
		// Apply styles line by line
		viewportLines := strings.Split(viewportContent, "\n")
		for i, line := range viewportLines {
			if i > 0 {
				content.WriteString("\n")
			}
			content.WriteString(styles.DetailContent.Render(line))
		}
		
		// Pad to fill the viewport height if needed
		linesRendered := len(viewportLines)
		for i := linesRendered; i < d.viewport.Height; i++ {
			content.WriteString("\n")
		}
	}

	// Add placeholders for future features (fixed at bottom)
	content.WriteString("\n\n")
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
	// Skip if same ticket is already selected
	if d.SelectedTicket != nil && ticket != nil && d.SelectedTicket.ID == ticket.ID {
		return
	}
	
	d.SelectedTicket = ticket
	d.SelectedProject = nil
	d.ViewType = "issues"
	
	// Update viewport content if ready
	if d.ready {
		d.viewport.SetContent(d.buildContent())
		d.viewport.GotoTop() // Reset scroll position for new content
	}
}

// SetSelectedProject sets the selected project for display
func (d *DetailPane) SetSelectedProject(project *mock.MockProject) {
	// Skip if same project is already selected
	if d.SelectedProject != nil && project != nil && d.SelectedProject.ID == project.ID {
		return
	}
	
	d.SelectedProject = project
	d.SelectedTicket = nil
	d.ViewType = "projects"
	
	// Update viewport content if ready
	if d.ready {
		d.viewport.SetContent(d.buildContent())
		d.viewport.GotoTop() // Reset scroll position for new content
	}
}

// SetFocus sets the focus state of the detail pane
func (d *DetailPane) SetFocus(focused bool) {
	d.Focused = focused
}

// SetDimensions sets the width and height of the detail pane
func (d *DetailPane) SetDimensions(width, height int) {
	// Skip if dimensions haven't changed
	if d.Width == width && d.Height == height && d.ready {
		return
	}
	
	d.Width = width
	d.Height = height
	
	// Initialize viewport if not ready and we have valid dimensions
	if !d.ready && width > 0 && height > 0 {
		// Calculate available height for viewport
		// Account for: border (2), title (2), placeholder section (7 lines + 2 spacing)
		availableHeight := height - 2 - 2 - 9
		if availableHeight < 1 {
			availableHeight = 1
		}
		
		// Account for border and padding in width
		availableWidth := width - 4
		if availableWidth < 1 {
			availableWidth = 1
		}
		
		d.viewport = viewport.New(availableWidth, availableHeight)
		d.viewport.YPosition = 0
		d.viewport.HighPerformanceRendering = false
		d.viewport.MouseWheelEnabled = false // Disable mouse to avoid conflicts
		
		// Add j/k key bindings to the viewport
		d.viewport.KeyMap.Down.SetKeys("down", "j")
		d.viewport.KeyMap.Up.SetKeys("up", "k")
		d.viewport.KeyMap.PageDown.SetKeys("pgdown", " ")
		d.viewport.KeyMap.PageUp.SetKeys("pgup")
		
		d.ready = true
		
		// Set initial content
		d.viewport.SetContent(d.buildContent())
	} else if d.ready {
		// Update viewport dimensions if already initialized
		d.updateViewportDimensions()
		// Re-set content to reflow for new width
		d.viewport.SetContent(d.buildContent())
	}
}

// getFocusedPane returns the current focused pane for styling
func (d *DetailPane) getFocusedPane() Pane {
	if d.Focused {
		return PaneDetail
	}
	return PaneMain // Default if not focused
}

// buildContent builds the content string for the viewport
func (d *DetailPane) buildContent() string {
	var content strings.Builder
	
	// Content based on what's selected
	if d.ViewType == "issues" && d.SelectedTicket != nil {
		d.renderTicketDetailsForViewport(&content)
	} else if d.ViewType == "projects" && d.SelectedProject != nil {
		d.renderProjectDetailsForViewport(&content)
	} else {
		content.WriteString("Select an item to view details")
	}
	
	
	return content.String()
}

// renderTicketDetailsForViewport renders ticket details for the viewport (no styling)
func (d *DetailPane) renderTicketDetailsForViewport(content *strings.Builder) {
	ticket := d.SelectedTicket
	
	// Title
	content.WriteString(fmt.Sprintf("Title: %s\n\n", ticket.Title))
	
	// Metadata
	content.WriteString(fmt.Sprintf("ID: %s\n", ticket.ID))
	content.WriteString(fmt.Sprintf("Status: %s\n", ticket.Status))
	content.WriteString(fmt.Sprintf("Priority: %s\n", ticket.Priority))
	content.WriteString(fmt.Sprintf("Assignee: %s\n", ticket.Assignee))
	content.WriteString(fmt.Sprintf("Created: %s\n\n", ticket.CreatedAt.Format("2006-01-02 15:04")))
	
	// Description
	content.WriteString("Description:\n")
	wrapWidth := d.Width - 8 // Default if viewport not ready
	if d.ready && d.viewport.Width > 0 {
		wrapWidth = d.viewport.Width
	}
	wrappedDesc := wordWrap(ticket.Description, wrapWidth)
	content.WriteString(wrappedDesc)
}

// renderProjectDetailsForViewport renders project details for the viewport (no styling)
func (d *DetailPane) renderProjectDetailsForViewport(content *strings.Builder) {
	project := d.SelectedProject
	
	// Title
	content.WriteString(fmt.Sprintf("Project: %s\n\n", project.Name))
	
	// Metadata
	content.WriteString(fmt.Sprintf("ID: %s\n", project.ID))
	content.WriteString(fmt.Sprintf("Status: %s\n", project.Status))
	content.WriteString(fmt.Sprintf("Progress: %.0f%%\n", project.Progress*100))
	content.WriteString(fmt.Sprintf("Created: %s\n\n", project.CreatedAt.Format("2006-01-02 15:04")))
	
	// Description
	content.WriteString("Description:\n")
	wrapWidth := d.Width - 8 // Default if viewport not ready
	if d.ready && d.viewport.Width > 0 {
		wrapWidth = d.viewport.Width
	}
	wrappedDesc := wordWrap(project.Description, wrapWidth)
	content.WriteString(wrappedDesc)
}

// updateViewportDimensions updates the viewport dimensions based on current pane size
func (d *DetailPane) updateViewportDimensions() {
	if !d.ready {
		return
	}
	
	// Calculate available height for viewport
	// Account for: border (2), title (2), placeholder section (7 lines + 2 spacing)
	availableHeight := d.Height - 2 - 2 - 9
	if availableHeight < 1 {
		availableHeight = 1
	}
	
	// Account for border and padding in width
	availableWidth := d.Width - 4
	if availableWidth < 1 {
		availableWidth = 1
	}
	
	d.viewport.Width = availableWidth
	d.viewport.Height = availableHeight
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
