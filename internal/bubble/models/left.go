package models

import (
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type LeftModel struct {
	Issues list.Model
	err    error
	dump   io.Writer
}

func (m *LeftModel) Init() tea.Cmd {
	return nil
}

func (m *LeftModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fprint(m.dump, "LEFT: ")
		spew.Fdump(m.dump, msg)
	}

	var cmd tea.Cmd
	m.Issues, cmd = m.Issues.Update(msg)
	return m, cmd
}

func (m *LeftModel) View() string {
	return docStyle.Render(m.Issues.View())
}

func (m *LeftModel) InitList() {
	m.Issues = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.Issues.Title = "To Do"
	m.Issues.SetItems([]list.Item{
		NewIssue("1", "Test", "Test", "Test", "Test", "Test", time.Now()),
		NewIssue("2", "Test", "Test", "Test", "Test", "Test", time.Now()),
		NewIssue("3", "Test", "Test", "Test", "Test", "Test", time.Now()),
		NewIssue("4", "Test", "Test", "Test", "Test", "Test", time.Now()),
		NewIssue("5", "Test", "Test", "Test", "Test", "Test", time.Now()),
		NewIssue("6", "Test", "Test", "Test", "Test", "Test", time.Now()),
		NewIssue("7", "Test", "Test", "Test", "Test", "Test", time.Now()),
		NewIssue("8", "Test", "Test", "Test", "Test", "Test", time.Now()),
	})
}
