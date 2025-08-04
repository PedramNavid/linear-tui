package models

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type RightModel struct {
	Details list.Model
	dump    io.Writer
}

func (m *RightModel) Init() tea.Cmd {
	return nil
}

func (m *RightModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fprint(m.dump, "RIGHT: ")
		spew.Fdump(m.dump, msg)
	}

	return m, nil
}

func (m *RightModel) View() string {
	return "Right Pane"
}
