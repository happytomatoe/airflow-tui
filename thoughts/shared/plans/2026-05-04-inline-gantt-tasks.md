# Inline Gantt Column for Tasks Implementation Plan

## Overview

Add an inline Gantt column to the task instances table showing duration bars directly in each row, matching the flows TUI behavior.

## Current State Analysis

- **Task table columns** (`ui/model.go:178-184`):
  - Task ID (width 40)
  - State (width 12)
  - Try # (width 6)
  - Duration (width 10)

- **updateTaskTable function** (`ui/model.go:1115-1141`): Builds table rows from taskInstances data
- **Existing Gantt chart** (`ui/gantt.go`): Has logic for calculating bar positions from StartDate/EndDate/Duration
- **TaskInstance fields** (`airflow/generated/client.go:1597-1667`):
  - StartDate *string
  - EndDate *string  
  - Duration *float32
  - State *TaskState

## Desired End State

Task table displays:
- Task ID | Duration | State | Tries | Gantt

Each row shows an inline progress bar representing task duration relative to max duration in the set.

### Key Discoveries:
- Gantt calculation logic exists in `ui/gantt.go:40-106` - can be reused
- TaskInstance has all needed fields: StartDate, EndDate, Duration
- Table column width calculation will need relative timing (find max duration, scale bars)

## What We're NOT Doing

- Removing the separate Gantt chart view (toggle with 'g' key) - keeping both
- Adding interactive Gantt features (clicking bars)
- Timeline-based Gantt (only duration bars, not absolute time axis)
- Sorting by Gantt position

## Implementation Approach

Add Gantt column to task table. Calculate bar by:
1. Find time window (min start, max end) across all tasks
2. For each task, render bar at start position with width based on duration
3. Use state-based coloring (green=success, red=failed, yellow=running, gray=queued)

## Phase 1: Add Gantt Column to Table Definition

### Overview
Add Gantt column to the task table columns definition.

### Changes Required:

#### 1. Task Table Column Definition
**File**: `ui/model.go`
**Changes**: Modify lines 178-184 to add Gantt column:

```go
taskTable := table.New(
    table.WithColumns([]table.Column{
        {Title: "Task ID", Width: 40},
        {Title: "Duration", Width: 10},
        {Title: "State", Width: 12},
        {Title: "Tries", Width: 6},
        {Title: "Gantt", Width: 20},
    }),
```

**Note**: Reordered columns to match flows (Duration before State), renamed "Try #" to "Tries".

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles: `go build .`
- [ ] No lint errors: `go fmt ./...`

#### Manual Verification:
- [ ] Table shows new Gantt column header
- [ ] Column ordering matches flows: Task ID, Duration, State, Tries, Gantt

---

## Phase 2: Implement Gantt Bar Rendering

### Overview
Create helper function to render Gantt bar string from task duration data.

### Changes Required:

#### 1. Gantt Bar Helper Function
**File**: `ui/model.go`
**Changes**: Add new function after line 1141 (after updateTaskTable):

```go
func renderGanttBar(task airflow.TaskInstance, maxDuration float32, width int) string {
    if task.Duration == nil || *task.Duration == 0 {
        return strings.Repeat(" ", width)
    }
    
    duration := *task.Duration
    if maxDuration <= 0 {
        maxDuration = 1
    }
    
    barWidth := int(float64(duration) / float64(maxDuration) * float64(width))
    if barWidth > width {
        barWidth = width
    }
    if barWidth < 1 {
        barWidth = 1
    }
    
    state := "none"
    if task.State != nil {
        state = string(*task.State)
    }
    
    bar := strings.Repeat("▃", barWidth)
    return ganttStateColor(state).Render(bar) + strings.Repeat(" ", width-barWidth)
}

func ganttStateColor(state string) lipgloss.Style {
    switch state {
    case "success":
        return lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
    case "failed":
        return lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
    case "running":
        return lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
    case "queued":
        return lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
    default:
        return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
    }
}
```

### Success Criteria:

#### Automated Verification:
- [x] Code compiles: `go build .`
- [x] No lint errors: `go fmt ./...`

#### Manual Verification:
- [ ] Function compiles without errors

---

## Phase 3: Update Task Table Rows with Gantt Column

### Overview
Modify updateTaskTable to include Gantt bar in each row.

### Changes Required:

#### 1. Update updateTaskTable Function
**File**: `ui/model.go`
**Changes**: Modify lines 1115-1141 to calculate max duration and add Gantt column:

```go
func (m *Model) updateTaskTable() {
    // Find max duration for Gantt scaling
    maxDuration := float32(0)
    for _, task := range m.taskInstances {
        if task.Duration != nil && *task.Duration > maxDuration {
            maxDuration = *task.Duration
        }
    }
    
    rows := make([]table.Row, 0, len(m.taskInstances))
    for _, task := range m.taskInstances {
        state := "-"
        if task.State != nil {
            state = string(*task.State)
        }
        tries := "-"
        if task.TryNumber != nil {
            tries = fmt.Sprintf("%d", *task.TryNumber)
        }
        duration := "-"
        if task.Duration != nil {
            duration = fmt.Sprintf("%.1fs", *task.Duration)
        }
        gantt := renderGanttBar(task, maxDuration, 15)
        rows = append(rows, table.Row{
            derefString(task.TaskId),
            duration,
            state,
            tries,
            gantt,
        })
    }
    m.taskTable.SetRows(rows)
    if len(rows) > 0 {
        m.taskTable.SetCursor(0)
    }
}
 ```

### Success Criteria:

#### Automated Verification:
- [x] Code compiles: `go build .`
- [x] No lint errors: `go fmt ./...`

#### Manual Verification:
- [ ] Tasks display with Gantt bars
- [ ] Bar width proportional to duration
- [ ] Colors match task state (green=success, red=failed, yellow=running)

---

## Phase 4: Adjust Column Widths and Test

### Overview
Fine-tune column widths for optimal display.

### Changes Required:

#### 1. Adjust Task Table Width
**File**: `ui/model.go`
**Changes**: May need to adjust individual column widths for better fit:

Current suggested widths:
- Task ID: 32 (reduced from 40 to make room)
- Duration: 8
- State: 10
- Tries: 6
- Gantt: 18

### Success Criteria:

#### Automated Verification:
- [x] Code compiles: `go build .`
- [x] `go fmt ./...` passes

#### Manual Verification:
- [ ] All columns visible without horizontal scroll on standard terminal (80+ cols)
- [ ] Gantt bars render correctly for tasks with varying durations
- [ ] State colors display correctly
- [ ] Toggle Gantt view ('g') still works independently

---

## Testing Strategy

### Manual Testing Steps:
1. Navigate to DAG → Runs → Tasks
2. Verify columns: Task ID | Duration | State | Tries | Gantt
3. Check that tasks with longer duration have wider bars
4. Verify color coding: success=green, failed=red, running=yellow, queued=gray
5. Press 'g' to toggle separate Gantt view - should still work
6. Test with tasks that have 0 or nil duration - should show empty bar

### Edge Cases:
- No tasks (empty state) - should show empty table gracefully
- Single task - bar should still render (maxDuration = task duration)
- Very long task names - truncated at column width
- Duration = 0 or nil - show empty/blank bar

## References

- Existing Gantt chart: `ui/gantt.go`
- TaskInstance model: `airflow/generated/client.go:1597-1667`
- Table implementation: `ui/model.go:178-192` (columns), `1115-1141` (rows)