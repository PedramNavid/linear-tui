package bubble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/linear-tui/linear-tui/internal/config"
	"github.com/linear-tui/linear-tui/internal/domain"
	"github.com/linear-tui/linear-tui/internal/ui/services"
)

type IssuesMsg struct {
	issues []domain.Issue
}

func (m IssuesMsg) Issues() []domain.Issue { return m.issues }

// GetLinearIssues is a command that fetches issues from Linear
func GetLinearIssues() tea.Msg {
	cfg, err := config.LoadConfig()
	if err != nil {
		return ErrMsg{err}
	}
	linearService, err := services.NewLinearService(cfg)
	if err != nil {
		return ErrMsg{err}
	}

	issues, err := linearService.GetTickets()
	if err != nil {
		return ErrMsg{err}
	}
	return IssuesMsg{issues}
}
