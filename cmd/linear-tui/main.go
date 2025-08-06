package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/linear-tui/linear-tui/internal/config"
	"github.com/linear-tui/linear-tui/internal/ui"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	// var dump *os.File
	// if cfg.DebugMode {
	// 	var err error
	// 	dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	// 	if err != nil {
	// 		os.Exit(1)
	// 	}
	// }
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
