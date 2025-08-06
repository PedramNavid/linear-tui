package detailpane

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui/messages"
)

type Styles struct {
	Border     lipgloss.Style
	EmptyState lipgloss.Style
	Content    lipgloss.Style
}

type Model struct {
	item     interface{}
	focused  bool
	viewport viewport.Model
	width    int
	height   int
	styles   Styles
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
		case "q", "esc":
			return m, func() tea.Msg {
				return messages.CloseDetailPaneMsg{}
			}
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.item == nil {
		return m.renderEmptyState()
	}

	content := m.renderItemDetails()
	m.viewport.SetContent(content)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		Render(m.viewport.View())
}

func (m *Model) SetItem(item interface{}) {
	m.item = item
	m.viewport.GotoTop()
}

func New() Model {
	vp := viewport.New(40, 20)
	return Model{
		viewport: vp,
		styles:   defaultStyles(),
	}
}

func defaultStyles() Styles {
	return Styles{
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(lipgloss.Color("#626262")),
		EmptyState: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")),
	}
}

func (m Model) renderEmptyState() string {
	return m.styles.EmptyState.Render("Select an item to view details")
}

func (m Model) renderItemDetails() string {
	if m.item == nil {
		return m.renderEmptyState()
	}
	// Stub implementation - format the item as string
	return m.styles.Content.Render(fmt.Sprintf("%v", m.item))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 2 // Account for border
	m.viewport.Height = height - 2
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}
