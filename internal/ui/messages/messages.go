package messages

// TabSwitchedMsg is sent when the user switches tabs
type TabSwitchedMsg struct{ Index int }

// ItemSelectedMsg is sent when an item is selected from the list
type ItemSelectedMsg struct{ Item interface{} }

// CloseDetailPaneMsg is sent when the detail pane should be closed
type CloseDetailPaneMsg struct{}

// DataLoadedMsg is sent when data is loaded from the API
type DataLoadedMsg struct{ Items []interface{} }

// ViewType represents the type of view being displayed
type ViewType int

const (
	IssueView ViewType = iota
	ProjectView
)