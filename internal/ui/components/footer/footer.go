package footer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width  int
	height int
	styles Styles
}

type Styles struct {
	Footer lipgloss.Style
}

func New() Model {
	return Model{
		height: 1,
		styles: defaultStyles(),
	}
}

func defaultStyles() Styles {
	return Styles{
		Footer: lipgloss.NewStyle().
			Background(lipgloss.Color("#303030")).
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m Model) View() string {
	help := "q: quit | tab: switch view | enter: select | esc: close detail"
	
	footerContent := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(help)
	
	return m.styles.Footer.Render(footerContent)
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
}