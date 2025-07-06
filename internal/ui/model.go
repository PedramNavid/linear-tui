package ui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/linear-tui/linear-tui/internal/ui/components"
)

// Model is the main application model that wraps the layout
type Model struct {
	Layout *components.Layout
}

// NewModel creates a new application model
func NewModel() *Model {
	return &Model{
		Layout: components.NewLayout(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.Layout.Init()
}

// Update handles all message updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	layout, cmd := m.Layout.Update(msg)
	m.Layout = layout
	return m, cmd
}

// View renders the application
func (m Model) View() string {
	return m.Layout.View()
}