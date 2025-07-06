package components

import (
	"fmt"
	"os"
	"strings"
	"time"
	
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Layout represents the main three-pane layout
type Layout struct {
	// Components
	Sidebar    *Sidebar
	MainPane   *MainPane
	DetailPane *DetailPane
	
	// State
	FocusedPane    Pane
	TerminalWidth  int
	TerminalHeight int
	
	// Styles
	Styles *Styles
}

// NewLayout creates a new layout with all components
func NewLayout() *Layout {
	layout := &Layout{
		Sidebar:        NewSidebar(),
		MainPane:       NewMainPane(),
		DetailPane:     NewDetailPane(),
		FocusedPane:    PaneMain, // Start with main pane focused
		Styles:         NewStyles(),
		TerminalWidth:  80,  // Default dimensions
		TerminalHeight: 24,
	}
	
	// Initialize focus states
	layout.updateFocusStates()
	
	return layout
}

// Init initializes the layout
func (l *Layout) Init() tea.Cmd {
	return nil
}

// Update handles all keyboard input and updates
func (l *Layout) Update(msg tea.Msg) (*Layout, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.handleResize(msg.Width, msg.Height)
		
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return l, tea.Quit
			
		case "tab":
			l.moveFocusForward()
			
		case "shift+tab":
			l.moveFocusBackward()
			
		case "esc":
			l.FocusedPane = PaneMain
			l.updateFocusStates()
			
		case "enter":
			if l.FocusedPane == PaneSidebar {
				l.handleSidebarSelection()
			}
			
		default:
			// Forward key messages to the focused pane
			l.updateFocusedPane(msg)
		}
	}

	// Update individual components
	sidebar, sidebarCmd := l.Sidebar.Update(msg)
	l.Sidebar = sidebar
	if sidebarCmd != nil {
		cmds = append(cmds, sidebarCmd)
	}

	mainPane, mainCmd := l.MainPane.Update(msg)
	l.MainPane = mainPane
	if mainCmd != nil {
		cmds = append(cmds, mainCmd)
	}

	detailPane, detailCmd := l.DetailPane.Update(msg)
	l.DetailPane = detailPane
	if detailCmd != nil {
		cmds = append(cmds, detailCmd)
	}

	// Update detail pane content based on main pane selection
	l.updateDetailPaneContent()

	return l, tea.Batch(cmds...)
}

// View renders the entire layout
func (l *Layout) View() string {
	// If no terminal size set yet, use reasonable defaults
	if l.TerminalWidth == 0 || l.TerminalHeight == 0 {
		l.TerminalWidth = 120
		l.TerminalHeight = 30
	}
	
	// Check minimum terminal size requirements
	minWidth := 80
	minHeight := 31
	
	if l.TerminalWidth < minWidth || l.TerminalHeight < minHeight {
		errorMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render("Terminal too small!")
		
		sizeMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Render(fmt.Sprintf("Current: %dx%d", l.TerminalWidth, l.TerminalHeight))
		
		reqMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Render(fmt.Sprintf("Minimum: %dx%d", minWidth, minHeight))
		
		return lipgloss.JoinVertical(lipgloss.Left, 
			"",
			errorMsg,
			"",
			sizeMsg,
			reqMsg,
			"",
			"Please resize your terminal window.",
		)
	}
	
	// Calculate layout dimensions
	dimensions := l.calculateLayout(l.TerminalWidth, l.TerminalHeight)
	
	// Update component dimensions
	l.Sidebar.SetDimensions(dimensions.SidebarWidth, dimensions.SidebarHeight)
	l.MainPane.SetDimensions(dimensions.MainWidth, dimensions.MainHeight)
	l.DetailPane.SetDimensions(dimensions.DetailWidth, dimensions.DetailHeight)

	// Render components
	sidebarView := l.Sidebar.View(l.Styles)
	mainView := l.MainPane.View(l.Styles)
	detailView := l.DetailPane.View(l.Styles)

	// Debug actual rendered dimensions
	l.debugRenderedSizes(sidebarView, mainView, detailView)

	// Create the right column (main + detail panes stacked vertically)
	rightColumn := lipgloss.JoinVertical(lipgloss.Top, mainView, detailView)
	
	// Join sidebar on the left with the right column
	layout := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, rightColumn)
	
	return layout
}

// handleResize updates the terminal dimensions and recalculates layout
func (l *Layout) handleResize(width, height int) {
	l.TerminalWidth = width
	l.TerminalHeight = height
}

// moveFocusForward moves focus to the next pane
func (l *Layout) moveFocusForward() {
	switch l.FocusedPane {
	case PaneSidebar:
		l.FocusedPane = PaneMain
	case PaneMain:
		l.FocusedPane = PaneDetail
	case PaneDetail:
		l.FocusedPane = PaneSidebar
	}
	l.updateFocusStates()
}

// moveFocusBackward moves focus to the previous pane
func (l *Layout) moveFocusBackward() {
	switch l.FocusedPane {
	case PaneSidebar:
		l.FocusedPane = PaneDetail
	case PaneMain:
		l.FocusedPane = PaneSidebar
	case PaneDetail:
		l.FocusedPane = PaneMain
	}
	l.updateFocusStates()
}

// updateFocusStates updates the focus state of all panes
func (l *Layout) updateFocusStates() {
	l.Sidebar.SetFocus(l.FocusedPane == PaneSidebar)
	l.MainPane.SetFocus(l.FocusedPane == PaneMain)
	l.DetailPane.SetFocus(l.FocusedPane == PaneDetail)
}

// handleSidebarSelection handles when an item is selected in the sidebar
func (l *Layout) handleSidebarSelection() {
	selectedKey := l.Sidebar.GetSelectedKey()
	l.MainPane.SetViewType(selectedKey)
	l.FocusedPane = PaneMain
	l.updateFocusStates()
}

// updateFocusedPane forwards input to the currently focused pane
func (l *Layout) updateFocusedPane(msg tea.KeyMsg) {
	switch l.FocusedPane {
	case PaneSidebar:
		// Sidebar handles its own input
	case PaneMain:
		// Main pane handles its own input
	case PaneDetail:
		// Detail pane handles its own input (currently read-only)
	}
}

// updateDetailPaneContent updates the detail pane based on main pane selection
func (l *Layout) updateDetailPaneContent() {
	if l.MainPane.ViewType == "issues" {
		selectedTicket := l.MainPane.GetSelectedTicket()
		l.DetailPane.SetSelectedTicket(selectedTicket)
	} else if l.MainPane.ViewType == "projects" {
		selectedProject := l.MainPane.GetSelectedProject()
		l.DetailPane.SetSelectedProject(selectedProject)
	}
}

// GetFocusedPane returns the currently focused pane
func (l *Layout) GetFocusedPane() Pane {
	return l.FocusedPane
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

// calculateLayout calculates the layout dimensions based on terminal size
func (l *Layout) calculateLayout(terminalWidth, terminalHeight int) LayoutDimensions {
	// Calculate sidebar dimensions (20% of width, full height)
	sidebarWidth := max(20, terminalWidth/5)
	sidebarHeight := terminalHeight
	
	// Calculate main and detail pane dimensions
	remainingWidth := terminalWidth - sidebarWidth
	mainWidth := remainingWidth
	detailWidth := remainingWidth
	
	// Ensure minimum heights for usability
	minMainHeight := 10
	minDetailHeight := 6  // Reduced from 8 to 6
	
	// Declare variables outside the conditional blocks
	var mainHeight, detailHeight int
	
	// For small terminals (≤30 rows), use much more aggressive limits
	if terminalHeight <= 30 {
		// Very small terminal: detail pane gets only 6-8 rows max
		detailHeight = min(8, max(minDetailHeight, terminalHeight/5))
		mainHeight = terminalHeight - detailHeight
		
		// Ensure main pane meets minimum
		if mainHeight < minMainHeight {
			mainHeight = minMainHeight
			detailHeight = terminalHeight - mainHeight
			if detailHeight < minDetailHeight {
				detailHeight = minDetailHeight
			}
		}
	} else {
		// Large terminal: use 70/30 split but cap detail pane at 15 rows
		mainHeight = (terminalHeight * 7) / 10  // 70% to main
		detailHeight = min(15, terminalHeight - mainHeight)
		
		// Recalculate main height after capping detail
		mainHeight = terminalHeight - detailHeight
	}
	
	// Debug logging
	l.debugLog(terminalWidth, terminalHeight, sidebarWidth, sidebarHeight, mainWidth, mainHeight, detailWidth, detailHeight, minMainHeight, minDetailHeight)
	
	return LayoutDimensions{
		SidebarWidth:  sidebarWidth,
		SidebarHeight: sidebarHeight,
		MainWidth:     mainWidth,
		MainHeight:    mainHeight,
		DetailWidth:   detailWidth,
		DetailHeight:  detailHeight,
	}
}

// debugLog writes layout calculations to log.txt
func (l *Layout) debugLog(terminalWidth, terminalHeight, sidebarWidth, sidebarHeight, mainWidth, mainHeight, detailWidth, detailHeight, minMainHeight, minDetailHeight int) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	
	timestamp := time.Now().Format("15:04:05.000")
	logEntry := fmt.Sprintf("[%s] LAYOUT DEBUG:\n", timestamp)
	logEntry += fmt.Sprintf("  Terminal: %dx%d\n", terminalWidth, terminalHeight)
	logEntry += fmt.Sprintf("  Sidebar: %dx%d\n", sidebarWidth, sidebarHeight)
	logEntry += fmt.Sprintf("  Main: %dx%d (min: %d)\n", mainWidth, mainHeight, minMainHeight)
	logEntry += fmt.Sprintf("  Detail: %dx%d (min: %d)\n", detailWidth, detailHeight, minDetailHeight)
	logEntry += fmt.Sprintf("  Total calculated height: %d (main+detail)\n", mainHeight+detailHeight)
	logEntry += fmt.Sprintf("  Height difference: %d\n", terminalHeight-(mainHeight+detailHeight))
	
	// Show percentage calculations
	mainPercent := float64(mainHeight) / float64(terminalHeight) * 100
	detailPercent := float64(detailHeight) / float64(terminalHeight) * 100
	logEntry += fmt.Sprintf("  Main takes %.1f%% of height\n", mainPercent)
	logEntry += fmt.Sprintf("  Detail takes %.1f%% of height\n", detailPercent)
	
	// Show if minimums are being enforced
	originalMainHeight := (terminalHeight * 3) / 5
	originalDetailHeight := terminalHeight - originalMainHeight
	if mainHeight != max(minMainHeight, originalMainHeight) || detailHeight != max(minDetailHeight, originalDetailHeight) {
		logEntry += fmt.Sprintf("  CONSTRAINT APPLIED: original would be main=%d, detail=%d\n", originalMainHeight, originalDetailHeight)
	}
	
	logEntry += "\n"
	file.WriteString(logEntry)
}

// debugRenderedSizes logs the actual rendered sizes of components
func (l *Layout) debugRenderedSizes(sidebarView, mainView, detailView string) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	
	timestamp := time.Now().Format("15:04:05.000")
	logEntry := fmt.Sprintf("[%s] RENDERED SIZES:\n", timestamp)
	
	// Count actual lines in each rendered component
	sidebarLines := len(strings.Split(sidebarView, "\n"))
	mainLines := len(strings.Split(mainView, "\n"))
	detailLines := len(strings.Split(detailView, "\n"))
	
	// Get the width of each component (first line)
	sidebarWidth := 0
	mainWidth := 0
	detailWidth := 0
	
	if sidebarLines > 0 {
		sidebarWidth = len(strings.Split(sidebarView, "\n")[0])
	}
	if mainLines > 0 {
		mainWidth = len(strings.Split(mainView, "\n")[0])
	}
	if detailLines > 0 {
		detailWidth = len(strings.Split(detailView, "\n")[0])
	}
	
	logEntry += fmt.Sprintf("  Sidebar rendered: %dx%d lines\n", sidebarWidth, sidebarLines)
	logEntry += fmt.Sprintf("  Main rendered: %dx%d lines\n", mainWidth, mainLines)
	logEntry += fmt.Sprintf("  Detail rendered: %dx%d lines\n", detailWidth, detailLines)
	logEntry += fmt.Sprintf("  Total rendered height: %d lines\n", mainLines+detailLines)
	
	logEntry += "\n"
	file.WriteString(logEntry)
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}