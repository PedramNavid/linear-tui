package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/domain"
	"github.com/linear-tui/linear-tui/internal/ui/components/detailpane"
	"github.com/linear-tui/linear-tui/internal/ui/components/footer"
	"github.com/linear-tui/linear-tui/internal/ui/components/listview"
	"github.com/linear-tui/linear-tui/internal/ui/components/tabs"
	"github.com/linear-tui/linear-tui/internal/ui/messages"
)

type FocusArea int

const (
	FocusMain FocusArea = iota
	FocusDetailPane
	FocusTabs
)

type Model struct {
	width  int
	height int

	// Child Componennts
	tabs       tabs.Model
	listView   listview.Model
	detailPane detailpane.Model
	footer     footer.Model

	// State
	currentView    messages.ViewType
	focusArea      FocusArea
	detailPaneOpen bool

	// Data
	issues   []domain.Issue
	projects []domain.Project

	styles Styles
}

func NewModel() Model {
	return Model{
		tabs:       tabs.New([]string{"Issues", "Projects"}),
		listView:   listview.New(),
		detailPane: detailpane.New(),
		footer:     footer.New(),

		currentView:    messages.IssueView,
		focusArea:      FocusMain,
		detailPaneOpen: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateComponentSizes()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Cycle through focus areas
			switch m.focusArea {
			case FocusTabs:
				m.focusArea = FocusMain
			case FocusMain:
				if m.detailPaneOpen {
					m.focusArea = FocusDetailPane
				} else {
					m.focusArea = FocusTabs
				}
			case FocusDetailPane:
				m.focusArea = FocusTabs
			}
		}

	case messages.TabSwitchedMsg:
		if msg.Index == 0 {
			m.currentView = messages.IssueView
		} else {
			m.currentView = messages.ProjectView
		}

	case messages.ItemSelectedMsg:
		m.detailPaneOpen = true
		m.detailPane.SetItem(msg.Item)
		m.focusArea = FocusDetailPane

	case messages.CloseDetailPaneMsg:
		m.detailPaneOpen = false
		m.focusArea = FocusMain
	}

	// Update child components based on focus
	switch m.focusArea {
	case FocusTabs:
		m.tabs.Focus()
		m.listView.Blur()
		m.detailPane.Blur()
	case FocusMain:
		m.tabs.Blur()
		m.listView.Focus()
		m.detailPane.Blur()
	case FocusDetailPane:
		m.tabs.Blur()
		m.listView.Blur()
		m.detailPane.Focus()
	}

	// Update child components
	var cmd tea.Cmd
	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	m.listView, cmd = m.listView.Update(msg)
	cmds = append(cmds, cmd)

	if m.detailPaneOpen {
		m.detailPane, cmd = m.detailPane.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateComponentSizes() {
	tabHeight := 1
	footerHeight := 1
	contentHeight := m.height - tabHeight - footerHeight

	// Update listview size
	if m.detailPaneOpen {
		listWidth := m.width * 2 / 3
		m.listView.SetSize(listWidth, contentHeight)
		m.detailPane.SetSize(m.width-listWidth, contentHeight)
	} else {
		m.listView.SetSize(m.width, contentHeight)
	}

	// Update footer width
	m.footer.SetWidth(m.width)
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	tabBar := m.tabs.View()
	mainContent := m.listView.View()

	var content string
	if m.detailPaneOpen {
		mainWidth := m.width * 2 / 3
		detailView := m.detailPane.View()
		content = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(mainWidth).Render(mainContent),
			lipgloss.NewStyle().Width(m.width-mainWidth).Render(detailView),
		)
	} else {
		content = mainContent
	}

	footer := m.footer.View()

	return lipgloss.JoinVertical(lipgloss.Left, tabBar, content, footer)

}
