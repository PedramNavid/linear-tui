# Ticket Creation Guidelines

This document outlines the standard process for creating tickets in Linear for the Linear TUI project.
The Linear Project for this repo is: https://linear.app/pdrm/project/linear-tui-0a4e9480ae4b
Project ID: a3169a31-fb9a-4c1d-8ae2-0d1826b1eba0
Team ID: 59aa1359-c0da-4709-a07f-49438d7bf6f9

## Ticket Creation Workflow

1. **Initial Requirements Gathering**
   - Understand the feature request or bug report from the user
   - Clarify any ambiguous requirements through questions
   - Identify the type of ticket (Feature, Bug, Task, Improvement)
   - Research any documentation, websites, etc. provided by the user.

2. **Scope Definition**
   - Write a clear, concise title (max 100 characters)
   - Define the problem statement or opportunity
   - Outline what is included in the scope
   - Explicitly state what is NOT included (out of scope)
   - Estimate complexity (Small, Medium, Large)

3. **Detailed Description Structure**
   ```markdown
   ## Summary
   [Brief overview of what needs to be done and why]

   ## Background
   [Context and current state that led to this ticket]

   ## Scope
   ### In Scope
   - [Specific deliverable 1]
   - [Specific deliverable 2]

   ### Out of Scope
   - [What will NOT be addressed]

   ## Technical Details
   [Any technical specifications, API endpoints, data structures]

   ## Dependencies
   - [Prerequisite ticket/task]
   - [Required resources or access]
   - [External dependencies]

   ## Acceptance Criteria
   - [ ] [Specific measurable criterion 1]
   - [ ] [Specific measurable criterion 2]
   - [ ] [All tests pass]
   - [ ] [Documentation updated]

   ## Test Cases
   1. **Test Case Name**
      - Steps to reproduce
      - Expected outcome
      - Edge cases to consider

   ## References
   - [Link to related documentation]
   - [Link to design mockups]
   - [Link to related tickets]
   ```

4. **Dependency Analysis**
   - Identify blocking dependencies
   - Link to parent/child tickets if applicable
   - Note any external service dependencies
   - List required permissions or access
   - Identify team members who need to be involved

5. **Acceptance Criteria Definition**
   - Write specific, measurable criteria
   - Include functional requirements
   - Include non-functional requirements (performance, security)
   - Define success metrics
   - Specify required documentation updates
   - Include testing requirements

6. **Test Case Specification**
   - Define happy path scenarios
   - Include edge cases and error scenarios
   - Specify performance test requirements if applicable
   - Include integration test scenarios
   - Define user acceptance test criteria

7. **Reference Documentation**
   - Link to relevant design documents
   - Include API documentation references
   - Link to related tickets or epics
   - Reference any architectural decisions
   - Include links to external resources

8. **Ticket Metadata**
   - Assign to appropriate team
   - Set correct project
   - Apply relevant labels
   - Set target cycle or milestone
   - Assign initial estimate (if using story points)

9. **Human Review**
   - Present the draft ticket to human collaborator
   - Review all sections for completeness
   - Verify acceptance criteria are testable
   - Confirm scope is appropriate
   - Wait for approval before creating in Linear

10. **Ticket Creation in Linear**
    - Use the Linear MCP tool to create the ticket
    - Set all metadata fields correctly
    - Link to related tickets
    - Add to appropriate project
    - Notify relevant stakeholders

## Best Practices

- **Be Specific**: Avoid vague descriptions. Use concrete examples.
- **Think Testing First**: Define how you'll know when the work is complete.
- **Consider the Reader**: Write for someone unfamiliar with the context.
- **Keep It Focused**: One ticket should address one concern.
- **Update Regularly**: Keep tickets updated as understanding evolves.

## Example Ticket

```markdown
Title: Implement Issues List View with Filtering

## Summary
Create a paginated list view for Linear issues with filtering capabilities by status, assignee, and priority.

## Background
Users need to view and filter their Linear issues within the TUI application. Currently, there's no way to browse issues without using the Linear web interface.

## Scope
### In Scope
- Fetch issues from Linear API
- Display issues in a scrollable list
- Filter by status (Todo, In Progress, Done)
- Filter by assignee
- Filter by priority
- Keyboard navigation

### Out of Scope
- Issue creation/editing
- Bulk operations
- Custom field filtering

## Technical Details
- Use Linear GraphQL API `issues` query
- Implement pagination with 20 items per page
- Cache results for 5 minutes
- Use bubbletea list component

## Dependencies
- Linear API client implementation (#123)
- Authentication setup (#124)

## Acceptance Criteria
- [ ] Issues display with title, status, assignee
- [ ] Filters can be combined (AND logic)
- [ ] List updates when filters change
- [ ] Keyboard shortcuts work (j/k for navigation)
- [ ] Loading states display during API calls
- [ ] Error handling for API failures
- [ ] Unit tests achieve 80% coverage

## Test Cases
1. **Basic List Display**
   - Launch app, navigate to issues
   - Verify 20 issues display
   - Check all fields render correctly

2. **Filter by Status**
   - Apply "In Progress" filter
   - Verify only in-progress issues show
   - Clear filter, verify all issues return

3. **API Error Handling**
   - Disconnect network
   - Attempt to load issues
   - Verify error message displays

## References
- [Linear API Docs](https://developers.linear.app/docs/graphql/working-with-the-graphql-api)
- [Bubbletea List Example](https://github.com/charmbracelet/bubbletea/tree/master/examples/list)
- Parent Epic: #100
```

## Notes for Agents

- Always gather complete requirements before creating a ticket
- Break large features into smaller, manageable tickets
- Ensure each ticket can be completed independently when possible
- Include enough detail that another developer could implement it
- Consider the end user's perspective in acceptance criteria
- Don't create tickets without human review and approval
