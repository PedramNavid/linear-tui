package bubble

import (
	"io"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/linear-tui/linear-tui/internal/bubble/components"
	"github.com/linear-tui/linear-tui/internal/ui"
)

type ErrMsg struct {
	err error
}

func (e ErrMsg) Error() string {
	return e.err.Error()
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	Tab              key.Binding
}

type focusState int

const (
	leftPane focusState = iota
	rightPane
)

type mainModel struct {
	onStartup bool

	Title  string
	Styles Styles
	err    error

	Layout     *components.Layout
	width      int
	height     int
	focusState focusState

	dump io.Writer

	MenuBar *components.MenuBar
	//MainPane   *components.MainPane
	//DetailPane *components.DetailPane
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}
func NewModel(dump io.Writer) *mainModel {
	return &mainModel{
		dump:   dump,
		Layout: components.NewLayout(ui.NewStyles()),
	}
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	windowSizeMsg, isWindowSizeMsg := msg.(tea.WindowSizeMsg)

	// Since this program is using the full size of the viewport we
	// need to wait until we've received the window dimensions before
	// we can initialize the viewport. The initial dimensions come in
	// quickly, though asynchronously, which is why we wait for them
	// here.
	if m.onStartup && !isWindowSizeMsg {
		return m, nil
	}

	if m.onStartup && isWindowSizeMsg {
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.width = windowSizeMsg.Width - h
		m.height = windowSizeMsg.Height - v

		m.onStartup = false

		return m, nil
	}

	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.dump != nil {
		spew.Fprint(m.dump, "MAIN: ")
		spew.Fdump(m.dump, msg)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}
		if s := msg.String(); s == "tab" {
			if m.focusState == leftPane {
				m.focusState = rightPane
			} else {
				m.focusState = leftPane
			}
		}

		return m, tea.WindowSize()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ErrMsg:
		m.err = msg
		return m, tea.Quit

	}

	// var routedModel tea.Model

	// if m.focusState == leftPane {
	// 	routedModel, cmd = m.left.Update(msg)
	// 	m.left = *routedModel.(*models.LeftModel)
	// 	cmds = append(cmds, cmd)
	// } else {
	// 	routedModel, cmd = m.right.Update(msg)
	// 	m.right = *routedModel.(*models.RightModel)
	// 	cmds = append(cmds, cmd)
	// }

	return m, tea.Batch(cmds...)
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) View() string {

	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	// Style the panes
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder(), false, false, true, false).
		Width(m.width).
		Render("header")

	headerHeight := lipgloss.Height(header) + 2

	leftBoxStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.height - headerHeight).
		Border(lipgloss.RoundedBorder())

	rightBoxStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height - headerHeight).
		Border(lipgloss.RoundedBorder())

	// Set focus style
	if m.focusState == leftPane {
		leftBoxStyle = leftBoxStyle.BorderForeground(lipgloss.Color("228"))
	} else {
		rightBoxStyle = rightBoxStyle.BorderForeground(lipgloss.Color("228"))
	}

	// leftView := leftBoxStyle.Render(m.left.View())
	// rightView := rightBoxStyle.Render(m.right.View())

	layoutView := m.Layout.View(m.width, m.height)
	return layoutView
}

type Styles struct {
	Title lipgloss.Style
}
