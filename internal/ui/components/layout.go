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
	FocusedPane Pane
	AppState    AppState
	LastError   error

	// Layout configuration
	Config *LayoutConfig

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
		MenuBar:       NewMenuBar(),
		MainPane:      NewMainPane(),
		DetailPane:    NewDetailPane(),
		Modal:         NewCreateTicketModal(),
		ErrorModal:    NewErrorModal(),
		LinearService: linearService,
		FocusedPane:   PaneMain, // Start with main pane focused
		AppState:      initialState,
		LastError:     serviceError,
		Config:        NewLayoutConfig(),
		Styles:        NewStyles(),
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
	// Get initial window size
	cmds := []tea.Cmd{
		// Request initial window size
		func() tea.Msg {
			// This will trigger a WindowSizeMsg on startup
			return tea.WindowSizeMsg{}
		},
	}

	// Start loading data immediately if we have a Linear service
	if l.LinearService != nil {
		cmds = append(cmds, l.loadDataAsync())
	} else {
		// If no Linear service, we can't proceed
		l.AppState = StateError
		l.LastError = fmt.Errorf("linear API key not configured - please set LINEAR_API_KEY environment variable")
	}

	return tea.Batch(cmds...)
}

// Update handles all keyboard input and updates
func (l *Layout) Update(msg tea.Msg) (*Layout, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.handleResize(msg.Width, msg.Height)
		return l, nil

	case LoadingMsg:
		l.AppState = StateLoading
		// Start the actual loading process
		return l, l.loadDataCmd()

	case DataLoadedMsg:
		l.AppState = StateReady
		l.MainPane.Issues = msg.Issues
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
			// Global hotkey to create new issue
			l.Modal.Show()

		case "e":
			// Edit selected issue when in issues view
			if l.MainPane.ViewType == "issues" {
				selectedIssue := l.MainPane.GetSelectedIssue()
				if selectedIssue != nil {
					l.Modal.ShowEdit(selectedIssue)
				}
			}

		case "tab":
			l.moveFocusForward()

		case "shift+tab":
			l.moveFocusBackward()

		case "esc":
			l.FocusedPane = PaneMain
			l.updateFocusStates()

		case "ctrl+d":
			// Toggle detail pane visibility
			l.Config.ToggleDetailPane()
			// If detail pane was hidden and we were focused on it, move focus to main
			if !l.Config.ShowDetailPane && l.FocusedPane == PaneDetail {
				l.FocusedPane = PaneMain
			}
			l.updateFocusStates()
			return l, nil

		case "r":
			// Manual refresh of data
			if l.LinearService != nil {
				return l, l.refreshDataCmd()
			}

		case "enter":
			if l.FocusedPane == PaneMenu {
				l.handleMenuSelection()
			}

		default:
			// Forward key messages to the focused pane
			l.updateFocusedPane(msg)
		}

	case CreateTicketResult:
		// Handle ticket creation/update result
		if msg.Success {
			l.Modal.SubmitMessage = msg.Message

			// If this was an update, refresh the issue data
			var refreshCmd tea.Cmd
			if msg.IsUpdate && msg.IssueID != "" {
				refreshCmd = l.refreshSingleIssueCmd(msg.IssueID)
			} else if !msg.IsUpdate {
				// For creates, refresh the entire list to show the new issue
				refreshCmd = l.refreshDataCmd()
			}

			// Auto-close modal after success
			closeCmd := tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return CloseModalMsg{}
			})

			return l, tea.Batch(refreshCmd, closeCmd)
		} else {
			l.Modal.ErrorMessage = "Failed to create ticket: " + msg.Message
			l.Modal.IsSubmitting = false
		}

	case CloseModalMsg:
		l.Modal.Hide()

	case SingleIssueRefreshedMsg:
		// Handle single issue refresh
		if msg.Error != nil {
			// Log error but don't show modal - the update was successful
			// This is just a refresh failure
			return l, nil
		}

		// Update the issue in the main pane's list
		if msg.Issue != nil {
			l.MainPane.UpdateSingleIssue(*msg.Issue)
			// Also update detail pane if this issue is selected
			if l.MainPane.GetSelectedIssue() != nil &&
				l.MainPane.GetSelectedIssue().LinearID == msg.Issue.LinearID {
				l.DetailPane.SetSelectedIssue(msg.Issue)
			}
		}
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
		// Load issues from Linear
		issues, err := l.LinearService.GetTickets()
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
			Issues:   issues,
			Projects: projects,
		}
	}
}

// refreshDataCmd refreshes data from Linear service
func (l *Layout) refreshDataCmd() tea.Cmd {
	// First, trigger loading state
	l.AppState = StateLoading

	return func() tea.Msg {
		// Refresh the service's cached data first
		if err := l.LinearService.RefreshData(); err != nil {
			return DataLoadErrorMsg{
				Error: fmt.Errorf("failed to refresh data: %w", err),
			}
		}

		// Then load the data
		return LoadingMsg{Message: "Refreshing data from Linear..."}
	}
}

// refreshSingleIssueCmd refreshes a single issue by ID
func (l *Layout) refreshSingleIssueCmd(issueID string) tea.Cmd {
	return func() tea.Msg {
		issue, err := l.LinearService.GetTicketByID(issueID)
		if err != nil {
			return SingleIssueRefreshedMsg{
				Issue: nil,
				Error: fmt.Errorf("failed to refresh issue: %w", err),
			}
		}

		return SingleIssueRefreshedMsg{
			Issue: issue,
			Error: nil,
		}
	}
}

// View renders the entire layout
func (l *Layout) View() string {
	// Wait for valid window size (tmux can report 0x0 initially)
	if l.Config.ScreenWidth == 0 || l.Config.ScreenHeight == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Padding(2, 0).
			Render("Initializing...")
	}

	// Check if layout can fit
	if canFit, reason := l.Config.CanFitLayout(); !canFit {
		errorMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render(reason)

		sizeMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Render(fmt.Sprintf("Current: %dx%d", l.Config.ScreenWidth, l.Config.ScreenHeight))

		minWidth := l.Config.MinPaneWidth
		if l.Config.ShowDetailPane {
			minWidth = l.Config.MinPaneWidth * 2
		}
		reqMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Render(fmt.Sprintf("Minimum: %dx%d", minWidth, l.Config.MenuBarHeight+l.Config.FooterHeight+10))

		return lipgloss.JoinVertical(lipgloss.Center,
			"",
			errorMsg,
			"",
			sizeMsg,
			reqMsg,
			"",
			"Please resize your terminal window.",
		)
	}

	// Update component dimensions using layout config
	l.MenuBar.SetDimensions(l.Config.ScreenWidth, l.Config.MenuBarHeight)
	l.MainPane.SetDimensions(l.Config.MainPaneWidth, l.Config.MainContentHeight)
	l.DetailPane.SetDimensions(l.Config.DetailPaneWidth, l.Config.MainContentHeight)
	l.Modal.SetDimensions(l.Config.ScreenWidth, l.Config.ScreenHeight)
	l.ErrorModal.SetDimensions(l.Config.ScreenWidth, l.Config.ScreenHeight)

	// Render components
	menuView := l.MenuBar.View(l.Styles)
	mainView := l.MainPane.View(l.Styles)
	detailView := l.DetailPane.View(l.Styles)

	// Create the bottom row (main + detail panes side-by-side)
	var bottomRow string
	if l.Config.ShowDetailPane {
		bottomRow = lipgloss.JoinHorizontal(lipgloss.Top, mainView, detailView)
	} else {
		bottomRow = mainView
	}

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
			Width(l.Config.ScreenWidth).
			Height(l.Config.ScreenHeight).
			Align(lipgloss.Center, lipgloss.Center)

		// Combine background and modal
		return overlayStyle.Render(modalView)
	}

	return layout
}

// handleResize updates the terminal dimensions and recalculates layout
func (l *Layout) handleResize(width, height int) {
	l.Config.UpdateDimensions(width, height)

	// If detail pane is visible, ensure viewports are properly sized
	if l.Config.ShowDetailPane && l.DetailPane != nil {
		// This will trigger viewport resize in the detail pane
		l.DetailPane.SetDimensions(l.Config.DetailPaneWidth, l.Config.MainContentHeight)
	}
}

// moveFocusForward moves focus to the next pane
func (l *Layout) moveFocusForward() {
	switch l.FocusedPane {
	case PaneMenu:
		l.FocusedPane = PaneMain
	case PaneMain:
		if l.Config.ShowDetailPane {
			l.FocusedPane = PaneDetail
		} else {
			l.FocusedPane = PaneMenu
		}
	case PaneDetail:
		l.FocusedPane = PaneMenu
	}
	l.updateFocusStates()
}

// moveFocusBackward moves focus to the previous pane
func (l *Layout) moveFocusBackward() {
	switch l.FocusedPane {
	case PaneMenu:
		if l.Config.ShowDetailPane {
			l.FocusedPane = PaneDetail
		} else {
			l.FocusedPane = PaneMain
		}
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
		selectedIssue := l.MainPane.GetSelectedIssue()
		l.DetailPane.SetSelectedIssue(selectedIssue)
	case "projects":
		selectedProject := l.MainPane.GetSelectedProject()
		l.DetailPane.SetSelectedProject(selectedProject)
	}
}

// GetFocusedPane returns the currently focused pane
func (l *Layout) GetFocusedPane() Pane {
	return l.FocusedPane
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
		Width(l.Config.ScreenWidth).
		Height(l.Config.ScreenHeight).
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
