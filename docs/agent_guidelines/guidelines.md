Bubble Tea programs are comprised of a model that describes the application state and three simple methods on that model:

Init, a function that returns an initial command for the application to run.
Update, a function that handles incoming events and updates the model accordingly.
View, a function that renders the UI based on the data in the model.


## Debugging

You can’t really log to stdout with Bubble Tea because your TUI is busy occupying that! You can, however, log to a file by including something like the following prior to starting your Bubble Tea program:

if len(os.Getenv("DEBUG")) > 0 {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

}

# General TUI Guidelines

  - Keep Update() and View() methods fast to avoid UI lag
  - Use pointer receivers on models judiciously
  - Build a tree of models for complex applications
  - Use lipgloss Height() and Width() methods for layout calculations
  - Messages may not be received in order
  - Use teatest for end-to-end testing
  - Log to debug.log file since stdout is occupied by the TUI

