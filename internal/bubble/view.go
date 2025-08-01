package bubble

import "github.com/charmbracelet/lipgloss"

var appStyle = lipgloss.NewStyle().Margin(1, 2)

// View renders the current state of the model as a string
func (m Model) View() string {
	return appStyle.Render(m.issues.View())

}
