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

1. Keep the event loop fast
Bubble Tea processes messages in an event loop:

func (p *Program) eventLoop(model Model, cmds chan Cmd) (Model, error) {
    for {
        select {
        case msg := <-p.msgs:
            // handle quit, window resize, etc
            // ...
            var cmd Cmd
            model, cmd = model.Update(msg) // run update
            cmds <- cmd                    // process command (if any)
            p.renderer.write(model.View()) // send view to renderer
        }
    }
}
A message is received from the channel and sent to the Update() method on your model. The returned command is sent to a channel, to be invoked in a go routine elsewhere. Your model’s View() method is then invoked before repeating the loop and processing the next message.

Therefore Bubble Tea can only process messages as fast as as your Update() and View() methods. You want these methods to be fast otherwise your program may experience lag, resulting in an unresponsive UI. If your program generates a lot of messages they can back up and the program may appear to stall: a user presses a key and nothing happens for an indetermine amount of time.

The key to writing a fast model is to offload expensive operations to a tea.Cmd:

4. Use receiver methods on your model judiciously
In Go, a method receiver can be passed as either a value or a pointer. When in doubt, one typically uses a pointer receiver, with a value receiver reserved for particular use cases.

It can throw Go programmers then that the documented Bubble Tea models all have value receivers. It may be due to the fact Bubble Tea is based on the Elm Architecture, which is a purely functional pattern, where functions cannot change their internal state, and in Go a method with a value receiver cannot modify its receiver.

However, you are free to set whatever receiver type you like. If you use a pointer receiver for your model and you make, say, a change to the model in Init() then that change is persisted:

type model struct {
	content string
}

func (m *model) Init() tea.Cmd {
	m.content = "initialized\n"
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string { return m.content }

func main() {
	p := tea.NewProgram(&model{content: "uninitalized"})
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
Returns:

initialized
However, don’t make the mistake of introducing a race condition by making changes outside of the event loop:

type model struct {
	content string
}

func (m *model) Init() tea.Cmd {
	go func() {
		m.content = "initialized\n"
	}()
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string { return m.content }

func main() {
	p := tea.NewProgram(&model{content: "uninitalized"})
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
If you repeatedly run this program you’ll find it returns initialized some of the time and sometimes uninitialized. In the latter case, the event loop has already called View() before the go routine sets the content to initialized (see the event loop code above).

Unless there is a good reason to do otherwise, stick to the normal message flow: any changes to the model should be made in Update() and returned immediately in the first return value. Straying from this course not only defies the natural order of Bubbletea, it also risks making it slower (see Keep the event loop fast).

5. Messages are not necessarily received in the order they are sent
In Go, if you have more than one go routine sending to a channel, the order in which the sends and receives occur is unspecified:


6. Build a tree of models
Any non-trivial Bubble Tea program outgrows a single model. There’s a good chance you’re using Charm’s bubbles, which are models in their own right, each with a Init(), Update(), and View(). You embed these models within your own model. The same applies to your own code: you may want to push your own components into separate models. The original model then becomes the “top-level” model, whose role becomes merely a message router and screen compositor, responsible for routing messages to the correct “child” models, and populating a layout with content from the child models’ View() methods.

And in turn the child models may embed models too, forming a tree of models: the root model receives all messages, which are relayed down the tree to the relevant child model, and the resulting model and command(s) are passed back up the tree, to be returned by the root model’s Update() method. The same traversal then occurs with the rendering: the root model’s View() method is called, which in turn calls child models’ View() methods, and the resulting strings are passed back up the tree to be joined together and returned to the renderer.

The root model maintains a list or map of child models. Depending on your program, you may nominate a child model to be the “current” model, which is the one that is currently visible and the one the user interacts with. You might maintain a stack of previously visited models: when the user presses a key your program pushes another model onto the stack, and the top of the stack is then the current model. When the user presses a key to go “back” the model is “popped” off the stack and the previous model becomes the current model.

You can opt to create your child models up front upon program startup. Or you could create them dynamically upon demand, which makes sense if conceptually they don’t exist at startup or they may number into the thousands. In the case of PUG, a LogMessage model is only created when the user “drills down” into an individual log message. If you choose the dynamic approach it makes sense to maintain a cache of models to avoid unnecessarily re-creating models.

The root model receives all messages. There are three main paths for routing decisions:

Handle the message directly in the root model. e.g. “global keys” such those mapping to quit, help, etc,.
Route the message to the current model (if you have one). e.g. all keys other than global keys, such as PageUp and PageDown to scroll up and down some content.
Route the message to all child models, e.g. tea.WindowSizeMsg, which contains the current terminal dimensions and all child models may want to use it to calculate heights and widths for rendering content.
None of the above is wrought in iron. It may not make sense for your particular program. However, Bubble Tea leaves architectural decisions to you and you’ll need to make conscious decisions on how to manage the complexity that inevitably occurs once your program reaches a certain size.

7. Layout arithmetic is error-prone
You are responsible for ensuring your program fits in the terminal. Its dimensions arrive in a tea.WindowSizeMsg message, which is sent shortly after the program starts, and whenever the terminal is resized. Your model records the dimensions and uses them when rendering to calculate the sizes of widgets.

In this app, there are three widgets: a header, content, and a footer:

type model struct {
	width, height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("header")
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("footer")

	content := lipgloss.NewStyle().
		Width(m.width).
        // accommodate header and footer
		Height(m.height-1-1).
		Align(lipgloss.Center, lipgloss.Center).
		Render("content")

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
Which produces:

working layout
The header and footer are of fixed sizes, and the content widget takes whatever space is leftover.

The code is then amended to add a border to the bottom of the header:

func (m model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render("header")
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("footer")

	content := lipgloss.NewStyle().
		Width(m.width).
        // accommodate header and footer
		Height(m.height-1-1).
		Align(lipgloss.Center, lipgloss.Center).
		Render("content")

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
But this breaks the layout, forcing the header off the terminal:

broken layout
The problem is that the arithmetic has not been updated to accommodate the border. The code is brittle, using hard coded heights which can easily be forgotten about when updating code. The fix is to use lipgloss’s Height() and Width() methods to reference heights and widths of widgets:

func (m model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render("header")
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("footer")

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-lipgloss.Height(header)-lipgloss.Height(footer)).
		Align(lipgloss.Center, lipgloss.Center).
		Render("content")

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
Which fixes the layout:

fixed layout
Now when changes are made to widget sizes the layout adapts accordingly.

As your program gets more complex, with more widgets and more models, it’s important to be disciplined with setting dimensions, to avoid frustratingly trying to track down what has caused the layout to break.


9. Use teatest for end-to-end tests
For end to end testing of your TUI, Charm have developed teatest, which they introduced last year in a blog article.

Here’s an program that runs and then quits upon confirmation from the user:

type model struct {
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, nil
		}
		if m.quitting {
			switch {
			case msg.String() == "y":
				return m, tea.Quit
			default:
				m.quitting = false
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Quit? (y/N)"
	} else {
		return "Running."
	}
}
And heres the test:

func TestQuit(t *testing.T) {
	m := model{}
	tm := teatest.NewTestModel(t, m)

	waitForString(t, tm, "Running.")

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	waitForString(t, tm, "Quit? (y/N)")

	tm.Type("y")

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func waitForString(t *testing.T, tm *teatest.TestModel, s string) {
	teatest.WaitFor(
		t,
		tm.Output(),
		func(b []byte) bool {
			return strings.Contains(string(b), s)
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*10),
	)
}
As you can see the test emulates the user pressing keys and checking that the program responds in kind, before checking the program has finished.

While this particular test is only checking for sub-strings, the blog article linked above shows how teatest supports using “golden files”, where the entire output is captured the first time the test is run, and subsequent tests then check the content matches the captured output. That’s useful for regression testing of content, but does mean you need to re-generate the golden files everytime you make even minor changes to the content of your program.

Note: as of writing teatest is part of Charm’s experimental repo, which means there is no promise of backwards compatibility.
