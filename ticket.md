# Add color styling for Backlog/In Progress/Done statuses in issues panel

## Summary
Currently, issue workflow statuses (Backlog, Todo, In Progress, Done) are displayed as plain text without any color styling in both the main issues list and detail pane. This ticket will add appropriate color themes to these statuses to improve visual clarity and help users quickly identify issue states.

## Background
The Linear TUI application already has a comprehensive styling system using lipgloss, with color styles for priority levels (High/Medium/Low). However, issue workflow statuses are rendered without any styling, making it harder to visually distinguish between different issue states at a glance.

## Scope
### In Scope
- Define new color styles for issue workflow statuses (Backlog, Todo, In Progress, Done)
- Create or extend GetStatusStyle function to handle issue workflow statuses
- Update main pane issue rendering to apply status colors
- Update detail pane status display to apply status colors
- Ensure color choices are appropriate and accessible

### Out of Scope
- Changing the existing priority color scheme
- Adding icons or other visual indicators beyond color
- Implementing theming or user-customizable colors
- Modifying project status colors

## Technical Details
- Extend the Styles struct in `internal/ui/styles.go` with new status styles:
  - StatusBacklog
  - StatusTodo
  - StatusInProgress
  - StatusDone (already exists but needs adjustment)
- Update GetStatusStyle function to handle issue workflow statuses
- Modify `internal/ui/components/mainpane.go` line ~162 to apply status styling
- Modify `internal/ui/components/detailpane.go` line ~141 to apply status styling
- Suggested color scheme:
  - Backlog: Gray (#808080)
  - Todo: Blue (#00BFFF)
  - In Progress: Yellow/Orange (#FFA500)
  - Done: Green (#00FF00) with strikethrough

## Dependencies
- None - all required infrastructure is already in place

## Acceptance Criteria
- [ ] Status colors are defined for all issue workflow states
- [ ] Issue statuses in the main pane list view are displayed with appropriate colors
- [ ] Issue status in the detail pane is displayed with the same color scheme
- [ ] Colors are visually distinct and accessible
- [ ] Existing priority coloring continues to work as before
- [ ] Code follows existing patterns and conventions
- [ ] All tests pass
- [ ] make lint, make fmt, and make build run successfully

## Test Cases
1. **Main Pane Status Display**
   - Navigate to issues view
   - Verify each status type displays with correct color:
     - Backlog items show in gray
     - Todo items show in blue
     - In Progress items show in orange
     - Done items show in green with strikethrough

2. **Detail Pane Status Display**
   - Select an issue with each status type
   - Verify the status field shows the same color as in the list view

3. **Priority Colors Still Work**
   - Verify High/Medium/Low priority indicators maintain their existing colors
   - Ensure no visual conflicts between status and priority colors

## References
- Current styling implementation: `internal/ui/styles.go`
- Main pane rendering: `internal/ui/components/mainpane.go:162`
- Detail pane rendering: `internal/ui/components/detailpane.go:141`
- Mock data with statuses: `internal/ui/mock/data.go`