package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui/mock"
)

const (
	// Layout constants
	borderWidth       = 2 // Border takes 1 char on each side
	horizontalPadding = 4 // Padding: 2 chars on each side (from border style)
	verticalPadding   = 2 // Padding: 1 line on top and bottom (from border style)

	// Fixed element heights
	titleHeight       = 2 // Title line + margin
	placeholderHeight = 7 // Placeholder section lines
	placeholderMargin = 2 // Spacing before placeholder section
)

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

	// Styles reference (set from parent)
	styles *Styles
}

func NewDetailPane() *DetailPane {
	return &DetailPane{
		Focused:  false,
		ViewType: "issues",
	}
}

// getContentWidth returns the available width for content after accounting for borders and padding
func (d *DetailPane) getContentWidth() int {
	contentWidth := d.Width - borderWidth - horizontalPadding
	if contentWidth < 1 {
		return 1
	}
	return contentWidth
}

// getViewportHeight returns the available height for the viewport after accounting for fixed elements
func (d *DetailPane) getViewportHeight() int {
	// Account for: border, title, placeholder section with margin
	totalFixedHeight := borderWidth + titleHeight + placeholderHeight + placeholderMargin
	viewportHeight := d.Height - totalFixedHeight
	if viewportHeight < 1 {
		return 1
	}
	return viewportHeight
}

func (d *DetailPane) Update(msg tea.Msg) (*DetailPane, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case tea.WindowSizeMsg:
		if d.ready && d.styles != nil {
			d.updateViewportDimensions()
			d.viewport.SetContent(d.buildContent(d.styles))
		}
	}

	if d.ready {
		// Only process viewport updates when focused for keyboard events
		if _, ok := msg.(tea.KeyMsg); ok {
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

func (d *DetailPane) View(styles *Styles) string {
	// Store styles reference for use in other methods
	d.styles = styles

	// Build title section
	titleSection := d.buildTitleSection(styles)

	// Build viewport section
	viewportSection := d.buildViewportSection(styles)

	// Build placeholder section
	placeholderSection := d.buildPlaceholderSection(styles)

	// Compose sections using JoinVertical
	content := lipgloss.JoinVertical(
		lipgloss.Top,
		titleSection,
		viewportSection,
		placeholderSection,
	)

	// Apply border and sizing
	borderStyle := styles.GetBorderStyle(PaneDetail, d.getFocusedPane())
	width := max(d.Width-borderWidth, 1)

	return borderStyle.Width(width).Render(content)
}

// buildTitleSection builds the title section with scroll indicators
func (d *DetailPane) buildTitleSection(styles *Styles) string {
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

	return styles.DetailTitle.Width(d.getContentWidth()).Render(title)
}

// buildViewportSection builds the viewport content section
func (d *DetailPane) buildViewportSection(styles *Styles) string {
	if !d.ready {
		return styles.Placeholder.
			Width(d.getContentWidth()).
			Height(d.getViewportHeight()).
			Render("Loading...")
	}

	// The viewport already contains styled content from buildContent
	viewportContent := d.viewport.View()

	// Ensure consistent height by padding if needed
	lines := strings.Split(viewportContent, "\n")
	for len(lines) < d.viewport.Height {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// buildPlaceholderSection builds the placeholder section for future features
func (d *DetailPane) buildPlaceholderSection(styles *Styles) string {
	placeholders := []string{
		"--- Placeholders for future features ---",
		"• Comments",
		"• Sub-issues",
		"• Attachments",
		"• Activity log",
	}

	var result strings.Builder
	for i, placeholder := range placeholders {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(styles.Placeholder.Width(d.getContentWidth()).Render(placeholder))
	}

	// Add margin before placeholders
	return "\n" + result.String()
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
	if d.ready && d.styles != nil {
		d.viewport.SetContent(d.buildContent(d.styles))
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
	if d.ready && d.styles != nil {
		d.viewport.SetContent(d.buildContent(d.styles))
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
		d.viewport = viewport.New(d.getContentWidth(), d.getViewportHeight())
		d.viewport.YPosition = 0
		d.viewport.MouseWheelEnabled = false // Disable mouse to avoid conflicts

		// Add j/k key bindings to the viewport
		d.viewport.KeyMap.Down.SetKeys("down", "j")
		d.viewport.KeyMap.Up.SetKeys("up", "k")
		d.viewport.KeyMap.PageDown.SetKeys("pgdown", " ")
		d.viewport.KeyMap.PageUp.SetKeys("pgup")

		d.ready = true

		// Set initial content if styles are available
		if d.styles != nil {
			d.viewport.SetContent(d.buildContent(d.styles))
		}
	} else if d.ready {
		// Update viewport dimensions if already initialized
		d.updateViewportDimensions()
		// Re-set content to reflow for new width
		if d.styles != nil {
			d.viewport.SetContent(d.buildContent(d.styles))
		}
	}
}

// getFocusedPane returns the current focused pane for styling
func (d *DetailPane) getFocusedPane() Pane {
	if d.Focused {
		return PaneDetail
	}
	return PaneMain // Default if not focused
}

// buildContent builds the styled content string for the viewport
func (d *DetailPane) buildContent(styles *Styles) string {
	var content strings.Builder

	// Content based on what's selected
	if d.ViewType == "issues" && d.SelectedTicket != nil {
		d.renderTicketDetailsStyled(&content, styles)
	} else if d.ViewType == "projects" && d.SelectedProject != nil {
		d.renderProjectDetailsStyled(&content, styles)
	} else {
		content.WriteString(styles.Placeholder.Render("Select an item to view details"))
	}

	return content.String()
}

// updateViewportDimensions updates the viewport dimensions based on current pane size
func (d *DetailPane) updateViewportDimensions() {
	if !d.ready {
		return
	}

	d.viewport.Width = d.getContentWidth()
	d.viewport.Height = d.getViewportHeight()
}

// renderTicketDetailsStyled renders ticket details with styles applied
func (d *DetailPane) renderTicketDetailsStyled(content *strings.Builder, styles *Styles) {
	ticket := d.SelectedTicket
	contentWidth := d.getContentWidth()

	// Title with proper width constraint
	titleStyle := styles.DetailContent.Width(contentWidth).MaxWidth(contentWidth)
	content.WriteString(titleStyle.Render(fmt.Sprintf("Title: %s", ticket.Title)))
	content.WriteString("\n\n")

	// Metadata with proper styling
	metaStyle := styles.DetailMeta.Width(contentWidth).MaxWidth(contentWidth)
	content.WriteString(metaStyle.Render(fmt.Sprintf("ID: %s", ticket.ID)))
	content.WriteString("\n")
	content.WriteString(metaStyle.Render(fmt.Sprintf("Status: %s", ticket.Status)))
	content.WriteString("\n")
	content.WriteString(metaStyle.Render(fmt.Sprintf("Priority: %s", ticket.Priority)))
	content.WriteString("\n")
	content.WriteString(metaStyle.Render(fmt.Sprintf("Assignee: %s", ticket.Assignee)))
	content.WriteString("\n")
	content.WriteString(metaStyle.Render(fmt.Sprintf("Created: %s", ticket.CreatedAt.Format("2006-01-02 15:04"))))
	content.WriteString("\n\n")

	// Description with lipgloss width handling
	content.WriteString(styles.DetailContent.Render("Description:"))
	content.WriteString("\n")

	// Use lipgloss to handle text wrapping
	descStyle := styles.DetailContent.Width(contentWidth).MaxWidth(contentWidth)
	content.WriteString(descStyle.Render(ticket.Description))
}

// renderProjectDetailsStyled renders project details with styles applied
func (d *DetailPane) renderProjectDetailsStyled(content *strings.Builder, styles *Styles) {
	project := d.SelectedProject
	contentWidth := d.getContentWidth()

	// Title with proper width constraint
	titleStyle := styles.DetailContent.Width(contentWidth).MaxWidth(contentWidth)
	content.WriteString(titleStyle.Render(fmt.Sprintf("Project: %s", project.Name)))
	content.WriteString("\n\n")

	// Metadata with proper styling
	metaStyle := styles.DetailMeta.Width(contentWidth).MaxWidth(contentWidth)
	content.WriteString(metaStyle.Render(fmt.Sprintf("ID: %s", project.ID)))
	content.WriteString("\n")
	content.WriteString(metaStyle.Render(fmt.Sprintf("Status: %s", project.Status)))
	content.WriteString("\n")
	content.WriteString(metaStyle.Render(fmt.Sprintf("Progress: %.0f%%", project.Progress*100)))
	content.WriteString("\n")
	content.WriteString(metaStyle.Render(fmt.Sprintf("Created: %s", project.CreatedAt.Format("2006-01-02 15:04"))))
	content.WriteString("\n\n")

	// Description with lipgloss width handling
	content.WriteString(styles.DetailContent.Render("Description:"))
	content.WriteString("\n")

	// Use lipgloss to handle text wrapping
	descStyle := styles.DetailContent.Width(contentWidth).MaxWidth(contentWidth)
	content.WriteString(descStyle.Render(project.Description))
}
