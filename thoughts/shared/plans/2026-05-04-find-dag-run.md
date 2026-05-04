# Find DAG Run Implementation Plan

## Overview

Add a "find run" feature on the DAG runs page that allows users to jump directly to a specific run by entering its Run ID. Activated via `f` key, it accepts a Run ID input and jumps to that row if found, or shows an error if not found.

## Current State Analysis

- **DAG runs view** (`ui/model.go:648-653`): Currently just renders `runsTable.View()` with no search/filter capability
- **DAGs page filter** (`ui/model.go:337-341`): Has existing search pattern using `/` key with `m.searching` bool and `m.input` textinput
- **Runs table structure** (`ui/model.go:149-163`): Columns are Run ID (width 40), State (10), Date (20), Type (12)
- **Run ID storage**: Stored in `run.DagRunId` field (`ui/model.go:815`)

## Desired End State

When viewing DAG runs (after selecting a DAG), pressing `f` opens an input prompt to enter a Run ID. On Enter:
- If Run ID matches exactly → cursor jumps to that row
- If Run ID not found → displays "Run not found: <run_id>" error message
- Pressing Escape cancels search and clears input

### Key Discoveries:
- `ui/model.go:254-272`: Existing pattern for handling search input (escape clears, enter applies)
- `ui/model.go:337`: Uses `m.input.Focus()` for input activation
- `ui/model.go:772-782`: `searchView()` renders input or filter status
- `ui/model.go:123`: Input placeholder is "find DAG by substring"

## What We're NOT Doing

- Fuzzy matching or partial matching (only exact Run ID match)
- Filtering to show subset of runs (jumping only)
- Searching by execution date, state, or other fields
- Reusing the `/` key (using `f` instead to avoid conflict)

## Implementation Approach

Use the existing `m.searching` and `m.input` infrastructure but with separate state tracking for runs page. The key is to reuse the existing textinput model rather than create a new one, but modify the behavior based on which tab is active.

## Phase 1: Add Key Handler for 'f' on Runs Page

### Overview
Add `f` key handling when on runs page to activate search input.

### Changes Required:

#### 1. Key Handler
**File**: `ui/model.go`
**Changes**: Add case for `f` key in the main key handling (around line 337, after the `/` case)

```go
case "f":
    if m.panel == dagRunPanel && m.client != nil {
        m.searching = true
        m.input.SetValue("")
        m.input.Placeholder = "find run by ID"
        return m, m.input.Focus()
    }
```

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles: `go build .`
- [ ] No lint errors: `go fmt ./...`

#### Manual Verification:
- [ ] On DAGs page, `f` does nothing
- [ ] On DAG runs page, pressing `f` shows input prompt

---

## Phase 2: Implement Run ID Search Logic

### Overview
When Enter is pressed during search on runs page, find the matching run by exact Run ID match and jump to it.

### Changes Required:

#### 1. Apply Filter for Runs
**File**: `ui/model.go`
**Changes**: Modify `applyFilter()` to handle runs page, or create new function `findRunByID()`

In the existing search handling (lines 254-272), add a condition to check if on runs panel:

```go
case "enter":
    if m.searching {
        m.filter = strings.TrimSpace(m.input.Value())
        if m.panel == dagRunPanel && m.filter != "" {
            // Find run by ID
            found := -1
            for i, run := range m.dagRuns {
                if derefString(run.DagRunId) == m.filter {
                    found = i
                    break
                }
            }
            if found >= 0 {
                m.runsTable.SetCursor(found)
            } else {
                m.err = fmt.Errorf("Run not found: %s", m.filter)
            }
        } else if m.panel == dagPanel {
            m.applyFilter()  // existing DAG filter
        }
        m.searching = false
        m.input.Blur()
        return m, nil
    }
```

**Note**: The existing code at line 262 already has `m.applyFilter()` for DAGs, need to add runs-specific logic there.

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles: `go build .`
- [ ] No lint errors: `go fmt ./...`

#### Manual Verification:
- [ ] Enter a valid Run ID → cursor jumps to that row
- [ ] Enter invalid Run ID → error message displayed
- [ ] Error clears when navigating away or pressing another key

---

## Phase 3: Add Error Display for Runs Search

### Overview
Ensure the "Run not found" error displays properly in the UI.

### Changes Required:

#### 1. Error Display Integration
**File**: `ui/model.go`
**Changes**: The existing error display at `ui/model.go:555-557` should automatically show the error. Verify that `m.err` is set correctly.

Key points:
- Line 556: `body = errorStyle.Render(fmt.Sprintf("Error: %v", m.err))`
- Error should appear above the table view

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles: `go build .`

#### Manual Verification:
- [ ] Invalid Run ID shows error message
- [ ] Error message is readable and explains what happened

---

## Phase 4: Cleanup and Edge Cases

### Overview
Handle edge cases and clean up the search state properly.

### Changes Required:

#### 1. Clear Error on Navigation
**File**: `ui/model.go`
**Changes**: Ensure error clears when:
- Navigating away from runs page (already handled at line 332 with `m.dagRuns = nil`)
- Pressing `f` again to search (should clear previous error)
- Pressing `r` to refresh (already clears at line 347 with `m.err = nil`)

#### 2. Update Help Text
**File**: `ui/help.go`
**Changes**: Add help text for `f` key on runs page

Around line 73 in help.go, add entry for runs page:
```go
case TabRuns:
    help = append(help, helpItem("f", "find run by ID"))
```

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles: `go build .`
- [ ] `go fmt ./...` passes

#### Manual Verification:
- [ ] Help text shows `f` for find run
- [ ] Error clears when navigating away

---

## Testing Strategy

### Manual Testing Steps:
1. Navigate to a DAG (press Enter on any DAG)
2. Navigate to Runs tab (press `3`)
3. Note a Run ID from the table
4. Press `f`, enter the Run ID, press Enter → should jump to that row
5. Press `f`, enter invalid ID like `nonexistent`, press Enter → should show error
6. Press Escape during search → should cancel and clear input
7. Press `r` to refresh → should clear error

### Edge Cases:
- Empty input (just press Enter) → should not error, just close input
- Very long Run ID → input has 128 char limit (line 168)
- Case sensitivity: Run IDs are exact match, so "run_1" ≠ "RUN_1"

## References

- Existing search pattern: `ui/model.go:254-272` (key handling with m.searching)
- DAG filter implementation: `ui/model.go:764-770` (applyFilter)
- Error display: `ui/model.go:555-557` (errorStyle usage)