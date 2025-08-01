package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/linear-tui/linear-tui/internal/config"
	"github.com/linear-tui/linear-tui/internal/domain"
	"github.com/linear-tui/linear-tui/internal/ui/services"
	"github.com/linear-tui/linear-tui/internal/ui/testdata"
)

type model struct {
	issues        []domain.Issue
	cursor        int
	selectedIssue map[int]struct{}
	err           error
}

func initialModel() model {
	return model{
		issues:        testdata.GetTestIssues(),
		cursor:        0,
		selectedIssue: make(map[int]struct{}),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.issues)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selectedIssue[m.cursor]
			if ok {
				delete(m.selectedIssue, m.cursor)
			} else {
				m.selectedIssue[m.cursor] = struct{}{}
			}
		}

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case issuesMsg:
		m.issues = msg.Issues()
		return m, nil

	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) Init() tea.Cmd {
	return getLinearIssues
}

func (m model) View() string {
	// The header
	s := "Linear Issues\n\n"

	// Iterate over our choices
	for i, issue := range m.issues {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selectedIssue[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, issue.Title)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func getLinearIssues() tea.Msg {
	cfg, err := config.LoadConfig()
	if err != nil {
		return errMsg{err}
	}
	linearService, err := services.NewLinearService(cfg)
	if err != nil {
		return errMsg{err}
	}

	issues, err := linearService.GetTickets()
	if err != nil {
		return errMsg{err}
	}
	return issuesMsg{issues}
}

type issuesMsg struct {
	issues []domain.Issue
}

func (m issuesMsg) Issues() []domain.Issue { return m.issues }

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
