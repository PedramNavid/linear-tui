package listview

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui/messages"
)

type Styles struct {
	EmptyState lipgloss.Style
	ListItem   lipgloss.Style
	Selected   lipgloss.Style
}

type Model struct {
	items    []string
	cursor   int
	focused  bool
	viewport viewport.Model
	table    table.Model
	viewType messages.ViewType
	styles   Styles
}

func New() Model {
	vp := viewport.New(80, 20)
	return Model{
		viewport: vp,
		items:    []string{},
		styles:   defaultStyles(),
	}
}

func defaultStyles() Styles {
	return Styles{
		EmptyState: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true),
		ListItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#874BFD")).
			Bold(true),
	}
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
		case "up", "k":
			m.cursor = (m.cursor - 1 + len(m.items)) % len(m.items)
		case "down", "j":
			m.cursor = (m.cursor + 1) % len(m.items)
		case "enter":
			if m.cursor < len(m.items) {
				return m, func() tea.Msg {
					return messages.ItemSelectedMsg{Item: m.items[m.cursor]}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if len(m.items) == 0 {
		return m.renderEmptyState()
	}

	return m.table.View()
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) renderEmptyState() string {
	return m.styles.EmptyState.Render("No items to display")
}

func (m *Model) SetSize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height
}
