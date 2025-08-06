package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui/messages"
)

type Styles struct {
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
}

type Model struct {
	tabs    []string
	active  int
	focused bool
	styles  Styles
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l", "right", "tab":
			m.active = (m.active + 1) % len(m.tabs)
			return m, func() tea.Msg {
				return messages.TabSwitchedMsg{Index: m.active}
			}
		case "h", "left", "shift+tab":
			m.active = (m.active - 1 + len(m.tabs)) % len(m.tabs)
			return m, func() tea.Msg {
				return messages.TabSwitchedMsg{Index: m.active}
			}

		}

	}

	return m, nil
}

func New(tabs []string) Model {
	return Model{
		tabs:   tabs,
		active: 0,
		styles: defaultStyles(),
	}
}

func defaultStyles() Styles {
	return Styles{
		TabActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#874BFD")).
			Padding(0, 2).
			Bold(true),
		TabInactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(0, 2),
	}
}

func (m Model) View() string {
	var tabs []string
	for i, tab := range m.tabs {
		if i == m.active {
			tabs = append(tabs, m.styles.TabActive.Render(tab))
		} else {
			tabs = append(tabs, m.styles.TabInactive.Render(tab))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}
