package components

// LayoutConfig centralizes all layout dimension calculations
// Following gh-dash patterns for consistent UI management
type LayoutConfig struct {
	// Terminal dimensions
	ScreenWidth  int
	ScreenHeight int

	// Fixed UI element heights
	MenuBarHeight int // Includes borders
	FooterHeight  int // Help text at bottom

	// Calculated dimensions
	MainContentHeight int // Available height for main content
	MainPaneWidth     int // Width of issues/list pane
	DetailPaneWidth   int // Width of detail/preview pane

	// Layout settings
	DetailPaneRatio float64 // Ratio of detail pane to total width (0.0-1.0)
	MinPaneWidth    int     // Minimum width for any pane
	ShowDetailPane  bool    // Whether detail pane is visible

	// Padding and borders
	BorderWidth int // Width consumed by borders (usually 1 per side)
	Padding     int // Internal padding for content
}

// NewLayoutConfig creates a layout config with sensible defaults
func NewLayoutConfig() *LayoutConfig {
	return &LayoutConfig{
		MenuBarHeight:   3,   // Title + border + padding
		FooterHeight:    1,   // Single line for help text
		DetailPaneRatio: 0.5, // 50/50 split by default
		MinPaneWidth:    30,  // Minimum readable width
		ShowDetailPane:  true,
		BorderWidth:     1,
		Padding:         1,
	}
}

// UpdateDimensions recalculates all derived dimensions based on screen size
func (lc *LayoutConfig) UpdateDimensions(width, height int) {
	lc.ScreenWidth = width
	lc.ScreenHeight = height

	// Calculate available content height
	lc.MainContentHeight = height - lc.MenuBarHeight - lc.FooterHeight
	if lc.MainContentHeight < 0 {
		lc.MainContentHeight = 0
	}

	// Calculate pane widths
	if lc.ShowDetailPane {
		// Calculate detail pane width based on ratio
		detailWidth := int(float64(width) * lc.DetailPaneRatio)

		// Ensure minimum widths
		if detailWidth < lc.MinPaneWidth {
			detailWidth = lc.MinPaneWidth
		}

		mainWidth := width - detailWidth
		if mainWidth < lc.MinPaneWidth {
			// If we can't fit both panes, prioritize main pane
			mainWidth = lc.MinPaneWidth
			detailWidth = width - mainWidth
			if detailWidth < 0 {
				// Can't fit both, hide detail pane
				lc.ShowDetailPane = false
				lc.MainPaneWidth = width
				lc.DetailPaneWidth = 0
				return
			}
		}

		lc.MainPaneWidth = mainWidth
		lc.DetailPaneWidth = detailWidth
	} else {
		lc.MainPaneWidth = width
		lc.DetailPaneWidth = 0
	}
}

// GetContentWidth returns the usable width for content after accounting for borders/padding
func (lc *LayoutConfig) GetContentWidth(paneWidth int) int {
	// Account for left/right borders and padding
	contentWidth := paneWidth - (2 * lc.BorderWidth) - (2 * lc.Padding)
	if contentWidth < 1 {
		return 1
	}
	return contentWidth
}

// GetContentHeight returns the usable height for content after accounting for borders/padding
func (lc *LayoutConfig) GetContentHeight() int {
	// Account for top/bottom borders and padding
	contentHeight := lc.MainContentHeight - (2 * lc.BorderWidth) - (2 * lc.Padding)
	if contentHeight < 1 {
		return 1
	}
	return contentHeight
}

// CanFitLayout checks if the current screen size can accommodate the layout
func (lc *LayoutConfig) CanFitLayout() (bool, string) {
	minWidth := lc.MinPaneWidth
	if lc.ShowDetailPane {
		minWidth = lc.MinPaneWidth * 2
	}

	minHeight := lc.MenuBarHeight + lc.FooterHeight + 10 // At least 10 lines for content

	if lc.ScreenWidth < minWidth {
		return false, "Terminal too narrow"
	}

	if lc.ScreenHeight < minHeight {
		return false, "Terminal too short"
	}

	return true, ""
}

// ToggleDetailPane toggles the detail pane visibility
func (lc *LayoutConfig) ToggleDetailPane() {
	lc.ShowDetailPane = !lc.ShowDetailPane
	lc.UpdateDimensions(lc.ScreenWidth, lc.ScreenHeight)
}
