# Rendered Template (t) and XCom (x) for Task Runs Implementation Plan

## Overview

Add hotkey `t` for viewing rendered templates and hotkey `x` for viewing XCom data from the task instances panel. This will provide developers quick access to task debugging information without leaving the TUI.

## Current State Analysis

What exists now:
- Task instances displayed in table with columns: Task ID, Duration, State, Tries, Gantt
- Hotkey `t` currently triggers DAG runs (will conflict with template viewing)
- Generated XCom API endpoints available but not integrated
- **Generated TaskInstance API with rendered_fields available but not integrated**

What's missing:
- XCom data loading and display functionality
- Template data loading and display functionality (GetTaskInstance already has rendered_fields)
- Hotkey handlers for task panel
- New UI views for XCom and template data

### Key Discoveries:
   - XCom API endpoints already exist in generated client: `GetXcomEntries` and `GetXcomEntry`
   - **Template rendering data is already available via `GetTaskInstance` API - returns `rendered_fields`**
   - Hotkey `t` is currently used for DAG triggering - will need to handle conflict
   - Task view hotkeys are handled in `ui/model.go:445` with switch statement
   - Footer help text in `ui/model.go:843` needs updating for new hotkeys
   - **Tested with `tutorial_taskflow_api` DAG: `rendered_fields` contains `op_args`, `op_kwargs`, `templates_dict`**

## What We're NOT Doing

- Implementing template editing functionality
- Adding XCom editing or management capabilities
- Creating new tabs for XCom/templates (integrating into existing task view)
- Implementing advanced XCom filtering/searching
- Adding template syntax highlighting
- Supporting custom XCom serializers/deserializers

## Implementation Approach

**Phase 1**: XCom Implementation (Higher Priority - API available)
**Phase 2**: Template Rendering (Lower Priority - Uses existing GetTaskInstance API)

## Phase 1: XCom Implementation

### Overview
Add XCom viewing functionality using existing generated API endpoints. This provides immediate value with minimal API development required.

### Changes Required:

#### 1. Client Interface Enhancement
**File**: `airflow/client.go`
**Changes**: Add XCom methods to Client interface

```go
// Add to Client interface
GetXcomEntries(ctx context.Context, dagID, dagRunID, taskID string, opts XComListOptions) ([]XCom, error)
GetXcomEntry(ctx context.Context, dagID, dagRunID, taskID, xcomKey string, opts XComGetOptions) (*XCom, error)
```

#### 2. Task View Hotkey Handling
**File**: `ui/model.go`
**Changes**: Add hotkey handler for 'x' around line 445

```go
case "x":
    return m.handleViewXComs()
```

#### 3. XCom Handler Methods
**File**: `ui/model.go`
**Changes**: Add new handler methods

```go
func (m *Model) handleViewXComs() (tea.Model, tea.Cmd) {
    if m.tabBar.Active() == TabTasks && m.nav.task != "" {
        m.loading = true
        m.xcomView = xcomLoading
        return m, m.loadXComs(m.nav.dag, m.nav.dagRun, m.nav.task)
    }
    return m, nil
}

func (m *Model) loadXComs(dagID, dagRunID, taskID string) tea.Cmd {
    return func() tea.Msg {
        xcoms, err := m.client.GetXcomEntries(context.Background(), dagID, dagRunID, taskID, XComListOptions{})
        return xcomLoadedMsg{dagID, dagRunID, taskID, xcoms, err}
    }
}
```

#### 4. XCom View Logic
**File**: `ui/model.go` 
**Changes**: Add XCom view method around line 800

```go
func (m *Model) xcomView() string {
    if m.xcomView == xcomNotLoaded {
        return m.theme.MutedStyle.Render("Press 'x' to view XCom data")
    }
    if m.xcomView == xcomLoading {
        return m.theme.MutedStyle.Render("Loading XCom data...")
    }
    // Display XCom data in table format
}
```

#### 5. Footer Help Text Update
**File**: `ui/model.go`
**Changes**: Update line 843

```go
case TabTasks:
    return m.theme.MutedStyle.Render("esc back  q quit  r refresh  1-5/←→ tabs  ? help  enter logs  c clear  m mark  o open  g gantt  t template  x xcom")
```

### Success Criteria:

#### Automated Verification:
- [ ] Client compiles without errors: `go build`
- ] XCom API integration works: `go test ./airflow`
- ] All existing tests pass: `just e2e-test`
- ] Hotkeys don't interfere with existing functionality

#### Manual Verification:
- [ ] XCom data loads correctly from task panel
- [ ] XCom data displays in readable format
- [ ] Navigation works (esc to return, arrow keys to scroll)
- [ ] Error handling for missing XCom data
- [ ] Performance acceptable for large XCom datasets

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Phase 2: Template Rendering Implementation

### Overview
Add template viewing functionality using existing `GetTaskInstance` API endpoint. The rendered_fields (op_args, op_kwargs, templates_dict) are already available - no additional API development required.

### Changes Required:

#### 1. Client Interface Enhancement
**File**: `airflow/client.go`
**Changes**: Add GetTaskInstance method to Client interface (already exists in generated client, needs wrapper)

```go
// Add to Client interface
GetTaskInstance(ctx context.Context, dagID, dagRunID, taskID string) (*TaskInstance, error)
```

#### 2. Template Hotkey Handling
**File**: `ui/model.go`
**Changes**: Handle hotkey conflict with context-aware behavior

```go
case "t":
    if m.tabBar.Active() == TabTasks && m.nav.task != "" {
        return m.handleViewTemplate()
    }
    // Keep existing DAG trigger behavior for other panels
    return m.handleTriggerDag()
```

#### 3. Template Handler Methods
**File**: `ui/model.go`
**Changes**: Add template handling methods

```go
func (m *Model) handleViewTemplate() (tea.Model, tea.Cmd) {
    if m.tabBar.Active() == TabTasks && m.nav.task != "" {
        m.loading = true
        m.templateView = templateLoading
        return m, m.loadTaskInstance(m.nav.dag, m.nav.dagRun, m.nav.task)
    }
    return m, nil
}

func (m *Model) loadTaskInstance(dagID, dagRunID, taskID string) tea.Cmd {
    return func() tea.Msg {
        ti, err := m.client.GetTaskInstance(context.Background(), dagID, dagRunID, taskID)
        return taskInstanceLoadedMsg{dagID, dagRunID, taskID, ti, err}
    }
}
```

#### 4. Template View Logic
**File**: `ui/model.go`
**Changes**: Add template view method around line 800

```go
func (m *Model) templateView() string {
    if m.templateView == templateNotLoaded {
        return m.theme.MutedStyle.Render("Press 't' to view template data")
    }
    if m.templateView == templateLoading {
        return m.theme.MutedStyle.Render("Loading template data...")
    }
    // Display rendered_fields: op_args, op_kwargs, templates_dict in table format
}
```

### Success Criteria:

#### Automated Verification:
- [ ] Client compiles without errors: `go build`
- [ ] Template API integration works: `go test ./airflow`
- [ ] All existing tests pass: `just e2e-test`
- [ ] Hotkeys don't interfere with existing functionality
- [ ] Rendered fields display correctly (op_args, op_kwargs, templates_dict)

#### Manual Verification:
- [ ] Template data loads correctly from task panel via GetTaskInstance API
- [ ] Rendered fields display in readable format (JSON/table)
- [ ] Navigation works (esc to return, arrow keys to scroll)
- [ ] Error handling for missing rendered fields
- [ ] Performance acceptable for large rendered field datasets

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Testing Strategy

### Unit Tests:
- XCom loading and error handling
- Template field listing and rendering
- Hotkey conflict resolution
- View state transitions

### Integration Tests:
- XCom data flow from API to UI
- Template rendering pipeline (GetTaskInstance → rendered_fields)
- Navigation between different views
- Error handling for API failures

### E2E Tests:
- `test_xcom_view_hotkey` - Press 'x' in task panel, verify XCom loads
- `test_template_view_hotkey` - Press 't' in task panel, verify rendered_fields display
- `test_hotkey_conflict_resolution` - Verify 't' triggers DAG in non-task panels
- `test_navigation_between_views` - Test esc to return from XCom/template views

### Manual Testing Steps:
1. Navigate to task instances panel
2. Press 'x' to view XCom data for a task
3. Verify XCom data loads and displays correctly
4. Test navigation within XCom view
5. Press 'esc' to return to task view
6. Press 't' to view rendered template data for a task
7. Verify rendered fields (op_args, op_kwargs, templates_dict) display correctly
8. Test error cases (no rendered fields, API errors)

## Performance Considerations

- XCom datasets can be large - implement pagination
- Use loading indicators for API calls
- Cache TaskInstance data to avoid repeated API calls when switching between XCom/template views
- Use efficient JSON display for rendered_fields data

## Migration Notes

- No data migration required - pure feature addition
- Backward compatibility maintained - existing hotkeys work as before
- No additional Airflow API development required - uses existing GetTaskInstance endpoint

## References

- Original TODO.md: "Add rendered template" and "Add xcom"
- Existing hotkey patterns in `ui/model.go:445`
- Footer help text in `ui/model.go:843`
- Generated XCom API in `airflow/generated/`
- Test patterns in `e2e/test_tui.py`
- **TaskInstance API**: `/dags/{dag_id}/dagRuns/{dag_run_id}/taskInstances/{task_id}` returns `rendered_fields`
- **Verified with**: `tutorial_taskflow_api` DAG transform task - contains `op_args`, `op_kwargs`, `templates_dict`