package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Pane represents which pane is currently focused
type Pane int

const (
	MenuPane Pane = iota
	MainPane
	DetailPane
)

// Styles contains all the styling for the application
type Styles struct {
	// Pane styles
	ActiveBorder   lipgloss.Style
	InactiveBorder lipgloss.Style

	// Menu bar styles
	MenuTitle    lipgloss.Style
	MenuItem     lipgloss.Style
	MenuSelected lipgloss.Style

	// Main pane styles
	MainTitle    lipgloss.Style
	ListItem     lipgloss.Style
	ListSelected lipgloss.Style

	// Detail pane styles
	DetailTitle   lipgloss.Style
	DetailContent lipgloss.Style
	DetailMeta    lipgloss.Style

	// General styles
	StatusHigh   lipgloss.Style
	StatusMedium lipgloss.Style
	StatusLow    lipgloss.Style
	StatusDone   lipgloss.Style

	// Placeholder styles
	Placeholder lipgloss.Style
}

// NewStyles creates a new Styles instance with default styling
func NewStyles() *Styles {
	return &Styles{
		// Pane borders
		ActiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2),

		InactiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#626262")).
			Padding(1, 2),

		// Menu bar styles
		MenuTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1).
			Bold(true),

		MenuItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(0, 1),

		MenuSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#874BFD")).
			Padding(0, 1).
			Bold(true),

		// Main pane styles
		MainTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1).
			Bold(true),

		ListItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(0, 1),

		ListSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#874BFD")).
			Padding(0, 1).
			Bold(true),

		// Detail pane styles
		DetailTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1).
			Bold(true),

		DetailContent: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1),

		DetailMeta: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(0, 1).
			Italic(true),

		// Status styles
		StatusHigh: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),

		StatusMedium: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true),

		StatusLow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),

		StatusDone: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Strikethrough(true),

		// Placeholder styles
		Placeholder: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true),
	}
}

// GetBorderStyle returns the appropriate border style based on focus state
func (s *Styles) GetBorderStyle(pane Pane, focusedPane Pane) lipgloss.Style {
	if pane == focusedPane {
		return s.ActiveBorder
	}
	return s.InactiveBorder
}

// GetStatusStyle returns the appropriate style for a status
func (s *Styles) GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "High":
		return s.StatusHigh
	case "Medium":
		return s.StatusMedium
	case "Low":
		return s.StatusLow
	case "Done", "Completed":
		return s.StatusDone
	default:
		return s.DetailMeta
	}
}

// LayoutDimensions calculates the dimensions for each pane
type LayoutDimensions struct {
	SidebarWidth  int
	SidebarHeight int
	MainWidth     int
	MainHeight    int
	DetailWidth   int
	DetailHeight  int
}

// CalculateLayout calculates the layout dimensions based on terminal size
func CalculateLayout(terminalWidth, terminalHeight int) LayoutDimensions {
	// Ensure minimum viable dimensions
	if terminalWidth < 80 {
		terminalWidth = 80
	}
	if terminalHeight < 24 {
		terminalHeight = 24
	}

	// Calculate sidebar dimensions (20% of width, full height)
	sidebarWidth := max(20, terminalWidth/5)
	sidebarHeight := terminalHeight

	// Calculate main and detail pane dimensions
	remainingWidth := terminalWidth - sidebarWidth
	mainWidth := remainingWidth
	detailWidth := remainingWidth

	// Main pane gets 60% of height, detail gets 40%
	mainHeight := max(10, (terminalHeight*3)/5)
	detailHeight := terminalHeight - mainHeight

	return LayoutDimensions{
		SidebarWidth:  sidebarWidth,
		SidebarHeight: sidebarHeight,
		MainWidth:     mainWidth,
		MainHeight:    mainHeight,
		DetailWidth:   detailWidth,
		DetailHeight:  detailHeight,
	}
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
