package components

import (
	"strings"

	"github.com/linear-tui/linear-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuBar struct {
	Items        []MenuItem
	SelectedItem int
	Width        int
	Height       int
	Focused      bool
}

type MenuItem struct {
	Title string
	Icon  string
	Key   string
}

func NewMenuBar() *MenuBar {
	return &MenuBar{
		Items: []MenuItem{
			{Title: "Issues", Icon: "🔍", Key: "issues"},
			{Title: "Projects", Icon: "📚", Key: "projects"},
		},
		SelectedItem: 0,
		Focused:      false,
	}
}

func (m *MenuBar) Update(msg tea.Msg) (*MenuBar, tea.Cmd) {
	if !m.Focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.SelectedItem > 0 {
				m.SelectedItem--
			}
		case "shift+tab":
			if m.SelectedItem < len(m.Items)-1 {
				m.SelectedItem++
			}
		}
	}

	return m, nil
}

func (m *MenuBar) View(styles *ui.Styles) string {
	var content strings.Builder

	for i, item := range m.Items {
		var style lipgloss.Style
		if i == m.SelectedItem && m.Focused {
			style = styles.MenuSelected
		} else {
			style = styles.MenuItem
		}

		itemText := style.Render(item.Title)
		content.WriteString(itemText)

		if i < len(m.Items)-1 {
			seperator := styles.MenuItem.Render(" | ")
			content.WriteString(seperator)
		}
	}
	borderStyle := styles.GetBorderStyle(ui.MenuPane, m.getFocusedPane())

	width := lipgloss.Width(content.String()) - 2
	if width < 1 {
		width = 1 // this is too low
	}

	menuContent := borderStyle.Width(width).Render(content.String())

	return menuContent
}

// Gets the key of the selected item
func (m *MenuBar) GetSelectedKey() string {
	if m.SelectedItem >= 0 && m.SelectedItem < len(m.Items) {
		return m.Items[m.SelectedItem].Key
	}
	return ""
}

func (m *MenuBar) SetFocused(focused bool) {
	m.Focused = focused
}

func (m *MenuBar) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

func (m *MenuBar) getFocusedPane() ui.Pane {
	if m.Focused {
		return ui.MenuPane
	}
	return ui.MainPane
}
