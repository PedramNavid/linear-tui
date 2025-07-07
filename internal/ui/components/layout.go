package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/config"
	"github.com/linear-tui/linear-tui/internal/ui/services"
)

// CloseModalMsg is sent to close the modal
type CloseModalMsg struct{}

// Layout represents the main layout with top menu and side-by-side panes
type Layout struct {
	// Components
	MenuBar    *MenuBar
	MainPane   *MainPane
	DetailPane *DetailPane
	Modal      *CreateTicketModal
	ErrorModal *ErrorModal

	// Services
	LinearService *services.LinearService

	// State
	FocusedPane    Pane
	TerminalWidth  int
	TerminalHeight int
	AppState       AppState
	LastError      error

	// Styles
	Styles *Styles
}

// NewLayout creates a new layout with all components
func NewLayout() *Layout {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fallback to default config if loading fails
		cfg = config.DefaultConfig()
	}

	// Try to initialize Linear service
	var linearService *services.LinearService
	var serviceError error
	if cfg.LinearAPIKey != "" {
		linearService, err = services.NewLinearService(cfg)
		if err != nil {
			serviceError = fmt.Errorf("failed to initialize Linear service: %w", err)
		}
	} else {
		serviceError = fmt.Errorf("linear API key not configured - please set LINEAR_API_KEY environment variable")
	}

	// Determine initial app state
	initialState := StateLoading
	if serviceError != nil {
		initialState = StateError
	}

	layout := &Layout{
		MenuBar:        NewMenuBar(),
		MainPane:       NewMainPane(),
		DetailPane:     NewDetailPane(),
		Modal:          NewCreateTicketModal(),
		ErrorModal:     NewErrorModal(),
		LinearService:  linearService,
		FocusedPane:    PaneMain, // Start with main pane focused
		AppState:       initialState,
		LastError:      serviceError,
		Styles:         NewStyles(),
		TerminalWidth:  80, // Default dimensions
		TerminalHeight: 24,
	}

	// Initialize focus states
	layout.updateFocusStates()

	// Set Linear service in modal for ticket creation
	if linearService != nil {
		layout.Modal.SetLinearService(linearService)
	}

	return layout
}

// Init initializes the layout
func (l *Layout) Init() tea.Cmd {
	// Start loading data immediately if we have a Linear service
	if l.LinearService != nil {
		return l.loadDataAsync()
	}

	// If no Linear service, we can't proceed
	l.AppState = StateError
	l.LastError = fmt.Errorf("linear API key not configured - please set LINEAR_API_KEY environment variable")
	return nil
}

// Update handles all keyboard input and updates
func (l *Layout) Update(msg tea.Msg) (*Layout, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.handleResize(msg.Width, msg.Height)

	case LoadingMsg:
		l.AppState = StateLoading
		// Start the actual loading process
		return l, l.loadDataCmd()

	case DataLoadedMsg:
		l.AppState = StateReady
		l.MainPane.Tickets = msg.Tickets
		l.MainPane.Projects = msg.Projects

	case DataLoadErrorMsg:
		l.AppState = StateError
		l.LastError = msg.Error
		l.ErrorModal.Show("Linear API Error", msg.Error.Error())

	case ErrorModalResult:
		switch msg.Action {
		case "retry":
			l.AppState = StateRetrying
			return l, l.loadDataAsync()
		case "quit":
			return l, tea.Quit
		}

	case tea.KeyMsg:
		// If error modal is visible, let it handle all input first
		if l.ErrorModal.IsVisible {
			errorModal, errorModalCmd := l.ErrorModal.Update(msg)
			l.ErrorModal = errorModal
			if errorModalCmd != nil {
				cmds = append(cmds, errorModalCmd)
			}
			// Return early to prevent background interaction
			return l, tea.Batch(cmds...)
		}

		// If create ticket modal is visible, let it handle all input first
		if l.Modal.IsVisible {
			modal, modalCmd := l.Modal.Update(msg)
			l.Modal = modal
			if modalCmd != nil {
				cmds = append(cmds, modalCmd)
			}
			// Return early to prevent background interaction
			return l, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return l, tea.Quit

		case "c":
			// Global hotkey to create new ticket
			l.Modal.Show()

		case "tab":
			l.moveFocusForward()

		case "shift+tab":
			l.moveFocusBackward()

		case "esc":
			l.FocusedPane = PaneMain
			l.updateFocusStates()

		case "enter":
			if l.FocusedPane == PaneMenu {
				l.handleMenuSelection()
			}

		default:
			// Forward key messages to the focused pane
			l.updateFocusedPane(msg)
		}

	case CreateTicketResult:
		// Handle ticket creation result
		if msg.Success {
			l.Modal.SubmitMessage = msg.Message
			// Auto-close modal after success
			return l, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return CloseModalMsg{}
			})
		} else {
			l.Modal.ErrorMessage = "Failed to create ticket: " + msg.Message
			l.Modal.IsSubmitting = false
		}

	case CloseModalMsg:
		l.Modal.Hide()
	}

	// Update individual components
	menuBar, menuBarCmd := l.MenuBar.Update(msg)
	l.MenuBar = menuBar
	if menuBarCmd != nil {
		cmds = append(cmds, menuBarCmd)
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

// loadDataAsync loads data asynchronously and sends messages about the result
func (l *Layout) loadDataAsync() tea.Cmd {
	return func() tea.Msg {
		// Send loading message
		return LoadingMsg{Message: "Loading data from Linear..."}
	}
}

// loadDataCmd performs the actual data loading
func (l *Layout) loadDataCmd() tea.Cmd {
	if l.LinearService == nil {
		return func() tea.Msg {
			return DataLoadErrorMsg{
				Error: fmt.Errorf("linear service not available"),
			}
		}
	}

	return func() tea.Msg {
		// Load tickets from Linear
		tickets, err := l.LinearService.GetTickets()
		if err != nil {
			return DataLoadErrorMsg{
				Error: fmt.Errorf("failed to load issues: %w", err),
			}
		}

		// Load projects from Linear
		projects, err := l.LinearService.GetProjects()
		if err != nil {
			return DataLoadErrorMsg{
				Error: fmt.Errorf("failed to load projects: %w", err),
			}
		}

		return DataLoadedMsg{
			Tickets:  tickets,
			Projects: projects,
		}
	}
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
	l.MenuBar.SetDimensions(dimensions.MenuWidth, dimensions.MenuHeight)
	l.MainPane.SetDimensions(dimensions.MainWidth, dimensions.MainHeight)
	l.DetailPane.SetDimensions(dimensions.DetailWidth, dimensions.DetailHeight)
	l.Modal.SetDimensions(l.TerminalWidth, l.TerminalHeight)
	l.ErrorModal.SetDimensions(l.TerminalWidth, l.TerminalHeight)

	// Render components
	menuView := l.MenuBar.View(l.Styles)
	mainView := l.MainPane.View(l.Styles)
	detailView := l.DetailPane.View(l.Styles)

	// Create the bottom row (main + detail panes side-by-side)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, mainView, detailView)

	// Stack menu bar on top, panes below
	layout := lipgloss.JoinVertical(lipgloss.Top, menuView, bottomRow)

	// Handle app states that override normal view
	switch l.AppState {
	case StateLoading, StateRetrying:
		loadingView := l.renderLoadingView()
		return loadingView
	case StateError:
		if l.ErrorModal.IsVisible {
			return l.ErrorModal.View(l.Styles)
		}
		// If error modal is not visible but we're in error state, show it
		l.ErrorModal.Show("Linear API Error", l.LastError.Error())
		return l.ErrorModal.View(l.Styles)
	}

	// Overlay create ticket modal if visible
	if l.Modal.IsVisible {
		modalView := l.Modal.View(l.Styles)

		// Create modal overlay with background
		overlayStyle := lipgloss.NewStyle().
			Width(l.TerminalWidth).
			Height(l.TerminalHeight).
			Align(lipgloss.Center, lipgloss.Center)

		// Combine background and modal
		return overlayStyle.Render(modalView)
	}

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
	case PaneMenu:
		l.FocusedPane = PaneMain
	case PaneMain:
		l.FocusedPane = PaneDetail
	case PaneDetail:
		l.FocusedPane = PaneMenu
	}
	l.updateFocusStates()
}

// moveFocusBackward moves focus to the previous pane
func (l *Layout) moveFocusBackward() {
	switch l.FocusedPane {
	case PaneMenu:
		l.FocusedPane = PaneDetail
	case PaneMain:
		l.FocusedPane = PaneMenu
	case PaneDetail:
		l.FocusedPane = PaneMain
	}
	l.updateFocusStates()
}

// updateFocusStates updates the focus state of all panes
func (l *Layout) updateFocusStates() {
	l.MenuBar.SetFocus(l.FocusedPane == PaneMenu)
	l.MainPane.SetFocus(l.FocusedPane == PaneMain)
	l.DetailPane.SetFocus(l.FocusedPane == PaneDetail)
}

// handleMenuSelection handles when an item is selected in the menu bar
func (l *Layout) handleMenuSelection() {
	selectedKey := l.MenuBar.GetSelectedKey()
	l.MainPane.SetViewType(selectedKey)
	l.FocusedPane = PaneMain
	l.updateFocusStates()
}

// updateFocusedPane forwards input to the currently focused pane
func (l *Layout) updateFocusedPane(msg tea.KeyMsg) {
	switch l.FocusedPane {
	case PaneMenu:
		// Menu bar handles its own input
	case PaneMain:
		// Main pane handles its own input
	case PaneDetail:
		// Detail pane handles its own input (currently read-only)
	}
}

// updateDetailPaneContent updates the detail pane based on main pane selection
func (l *Layout) updateDetailPaneContent() {
	switch l.MainPane.ViewType {
	case "issues":
		selectedTicket := l.MainPane.GetSelectedTicket()
		l.DetailPane.SetSelectedTicket(selectedTicket)
	case "projects":
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
	MenuWidth    int
	MenuHeight   int
	MainWidth    int
	MainHeight   int
	DetailWidth  int
	DetailHeight int
}

// calculateLayout calculates the layout dimensions based on terminal size
func (l *Layout) calculateLayout(terminalWidth, terminalHeight int) LayoutDimensions {
	// Calculate menu bar dimensions (full width, 3 lines height)
	menuWidth := terminalWidth
	menuHeight := 3

	// Calculate remaining space for main and detail panes
	remainingHeight := terminalHeight - menuHeight

	// Calculate main and detail pane dimensions (side by side)
	mainWidth := terminalWidth / 2           // 50% width for main pane
	detailWidth := terminalWidth - mainWidth // 50% width for detail pane

	// Both panes use full remaining height
	mainHeight := remainingHeight
	detailHeight := remainingHeight

	// Ensure minimum dimensions
	minMainWidth := 40
	minDetailWidth := 30

	if mainWidth < minMainWidth {
		mainWidth = minMainWidth
		detailWidth = terminalWidth - mainWidth
	}
	if detailWidth < minDetailWidth {
		detailWidth = minDetailWidth
		mainWidth = terminalWidth - detailWidth
	}

	return LayoutDimensions{
		MenuWidth:    menuWidth,
		MenuHeight:   menuHeight,
		MainWidth:    mainWidth,
		MainHeight:   mainHeight,
		DetailWidth:  detailWidth,
		DetailHeight: detailHeight,
	}
}

// renderLoadingView renders the loading screen
func (l *Layout) renderLoadingView() string {
	var message string
	switch l.AppState {
	case StateLoading:
		message = "Loading data from Linear..."
	case StateRetrying:
		message = "Retrying connection to Linear..."
	default:
		message = "Loading..."
	}

	loadingStyle := lipgloss.NewStyle().
		Width(l.TerminalWidth).
		Height(l.TerminalHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#CCCCCC"))

	spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏" // Simple spinner characters

	content := fmt.Sprintf("%s\n\n%c %s",
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Render("Linear TUI"),
		rune(spinner[0]), // Just use first character for now
		message)

	return loadingStyle.Render(content)
}

