package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/ui"
)

type Layout struct {
	MenuBar  *MenuBar
	MainPane *MainPane
	Styles   *ui.Styles
}

func NewLayout(styles *ui.Styles) *Layout {
	return &Layout{
		MenuBar:  NewMenuBar(),
		MainPane: NewMainPane(),
		Styles:   styles,
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

	mainPane, mainPaneCmd := l.MainPane.Update(msg)
	l.MainPane = mainPane
	if mainPaneCmd != nil {
		cmds = append(cmds, mainPaneCmd)
	}

	return l, tea.Batch(cmds...)
}

func (l *Layout) View(width, height int) string {
	l.MenuBar.SetDimensions(width-2, height-2)
	menuHeight := lipgloss.Height(l.MenuBar.View(l.Styles))
	l.MainPane.SetDimensions(width-2, height-menuHeight-2)
	menuView := l.MenuBar.View(l.Styles)
	mainPaneView := l.MainPane.View(l.Styles)

	return lipgloss.JoinVertical(lipgloss.Top, menuView, mainPaneView)
}
