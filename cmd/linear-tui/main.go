package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/linear-tui/linear-tui/internal/ui"
)

func main() {
	// Create the model
	model := ui.NewModel()
	
	// Create the program without alt screen to ensure top alignment
	p := tea.NewProgram(model)
	
	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}