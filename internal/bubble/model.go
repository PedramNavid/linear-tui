package bubble

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ErrMsg struct {
	err error
}

func (e ErrMsg) Error() string {
	return e.err.Error()
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	Tab              key.Binding
}

type focusState int

const (
	leftPane focusState = iota
	rightPane
)

type leftModel struct {
	Issues list.Model
	err    error
}

type rightModel struct {
	Details list.Model
}

type mainModel struct {
	Title        string
	Styles       Styles
	err          error
	left         leftModel
	right        rightModel
	width        int
	height       int
	focusState   focusState
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func (m *leftModel) Init() tea.Cmd {
	return nil
}

func (m *leftModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.Issues, cmd = m.Issues.Update(msg)
	return m, cmd
}

func (m *leftModel) View() string {
	return docStyle.Render(m.Issues.View())
}

func (m *leftModel) initList(width, height int) {
	m.Issues = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.Issues.Title = "To Do"
	m.Issues.SetItems([]list.Item{
		NewItem("1", "Test", "Test", "Test", "Test", "Test", time.Now()),
	})
}

func (m *rightModel) Init() tea.Cmd {
	return nil
}

func (m *rightModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *rightModel) View() string {
	return "Right Pane"
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}
func NewModel() *mainModel {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	return &mainModel{
		left: leftModel{
			Issues: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		},
		right:        rightModel{},
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.focusState == leftPane {
				m.focusState = rightPane
			} else {
				m.focusState = leftPane
			}
			return m, nil
		}

	case ErrMsg:
		m.err = msg
		return m, tea.Quit
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v
		m.left.Update(msg)
		return m, nil
	}

	var routedModel tea.Model
	if m.focusState == leftPane {
		routedModel, cmd = m.left.Update(msg)
		m.left = *routedModel.(*leftModel)
		cmds = append(cmds, cmd)
	} else {
		routedModel, cmd = m.right.Update(msg)
		m.right = *routedModel.(*rightModel)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func getTestData() []list.Item {
	return []list.Item{
		NewItem(
			"LIN-001",
			"Implement user authentication",
			"Add OAuth2 authentication flow for user login. This should include Google and GitHub providers with proper session management and token refresh capabilities.",
			"In Progress",
			"High",
			"Alice Johnson",
			time.Now().Add(-48*time.Hour),
		),
		NewItem(
			"LIN-002",
			"Fix database connection pooling",
			"The current database connection pool is not properly handling timeouts and connection limits. This causes the application to hang under high load.",
			"Todo",
			"High",
			"Bob Smith",
			time.Now().Add(-24*time.Hour),
		),
		NewItem(
			"LIN-003",
			"Add dark mode support",
			"Implement dark mode theme throughout the application with proper color scheme and user preference storage.",
			"Done",
			"Medium",
			"Carol Davis",
			time.Now().Add(-72*time.Hour),
		),
		NewItem(
			"LIN-004",
			"Optimize API response times",
			"Current API response times are averaging 500ms. Need to implement caching strategies and query optimization to bring this down to under 200ms.",
			"In Progress",
			"High",
			"David Wilson",
			time.Now().Add(-12*time.Hour),
		),
		NewItem(
			"LIN-005",
			"Create user onboarding flow",
			"Design and implement a comprehensive onboarding experience for new users including tutorials and tooltips.",
			"Todo",
			"Medium",
			"Unassigned",
			time.Now().Add(-96*time.Hour),
		),
	}
}

func (m mainModel) View() string {
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	// Style the panes
	leftBoxStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.height).
		Border(lipgloss.RoundedBorder())

	rightBoxStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height).
		Border(lipgloss.RoundedBorder())

	// Set focus style
	if m.focusState == leftPane {
		leftBoxStyle = leftBoxStyle.BorderForeground(lipgloss.Color("228"))
	} else {
		rightBoxStyle = rightBoxStyle.BorderForeground(lipgloss.Color("228"))
	}

	leftView := leftBoxStyle.Render(m.left.View())
	rightView := rightBoxStyle.Render(m.right.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)
}

type Styles struct {
	Title lipgloss.Style
}
