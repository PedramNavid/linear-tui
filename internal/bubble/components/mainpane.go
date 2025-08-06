package components

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/linear-tui/linear-tui/internal/bubble/models"
	"github.com/linear-tui/linear-tui/internal/ui"
)

type MainPane struct {
	Issues  list.Model
	Width   int
	Height  int
	Focused bool
}

func NewMainPane() *MainPane {
	return &MainPane{
		Issues: list.New([]list.Item{
			models.NewIssue("1", "Test", "Test", "Test", "Test", "Test", time.Now()),
			models.NewIssue("2", "Test", "Test", "Test", "Test", "Test", time.Now()),
		}, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m *MainPane) Update(msg tea.Msg) (*MainPane, tea.Cmd) {
	return m, nil
}

func (m *MainPane) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

func (m *MainPane) View(styles *ui.Styles) string {
	mainStyle := styles.GetBorderStyle(ui.MainPane, m.getFocusedPane())
	mainStyle.Width(m.Width).Height(m.Height)

	return mainStyle.Render(m.Issues.View())
}

func (m *MainPane) getFocusedPane() ui.Pane {
	if m.Focused {
		return ui.MainPane
	}
	return ui.MainPane
}
