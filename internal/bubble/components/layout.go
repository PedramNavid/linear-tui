package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui"
)

type Layout struct {
	MenuBar *MenuBar
	Styles  *ui.Styles
}

func NewLayout(styles *ui.Styles) *Layout {
	return &Layout{
		MenuBar: NewMenuBar(),
		Styles:  styles,
	}
}

func (l *Layout) Init() tea.Cmd {
	return nil
}

func (l *Layout) Update(msg tea.Msg) (*Layout, tea.Cmd) {
	var cmds []tea.Cmd
	menubar, menuBarCmd := l.MenuBar.Update(msg)
	l.MenuBar = menubar
	if menuBarCmd != nil {
		cmds = append(cmds, menuBarCmd)
	}

	return l, tea.Batch(cmds...)
}

func (l *Layout) View() string {
	//	menuBar := l.MenuBar.View(l.Styles)

	l.MenuBar.SetDimensions(150, 15)
	menuView := l.MenuBar.View(l.Styles)

	return lipgloss.JoinVertical(lipgloss.Top, menuView)
}
