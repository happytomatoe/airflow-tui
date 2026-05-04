# Fix Task Logs Display Implementation Plan

## Overview
This plan addresses the bug where task logs are incorrectly displayed in the Go TUI. The core issue seems to stem from how log content is fetched, parsed, and rendered. We will investigate the `GetTaskLog` implementation in `airflow/generated_client.go` and its interaction with the `ui/model.go` to correctly display task logs.

## Current State Analysis
- **`ui/model.go`**: Manages the `logViewport` and `logContent`, handles `logLoadedMsg` to update the log display.
- **`airflow/client.go`**: Defines the `Client` interface, which includes the `GetTaskLog` method.
- **`airflow/generated_client.go`**: Implements `GetTaskLog`, calls `parseV1LogContent` and attempts to skip a hostname line. This is the most likely place for the parsing issue.

### Key Discoveries:
- The `grep` for "task logs" returned no direct matches, indicating the display logic is likely distributed across `ui/model.go` and `airflow/generated_client.go`.
- The `grep` for "bubbletea" and "model|view|update" confirms `ui/model.go` as the central TUI component.
- The `grep` for "flowrs" shows many references within `flowrs/` subdirectories, which appears to be a separate (Rust-based) project. We should focus on the `airflow-tui` Go modules.

## Desired End State
Task logs are displayed correctly in the TUI, without any parsing errors, missing content, or extraneous information.

## What We're NOT Doing
- We are not refactoring the entire TUI or the Airflow API client.
- We are not introducing new features, only fixing the existing bug.

## Implementation Approach
We will proceed in phases, focusing on understanding the current log handling, debugging the parsing logic, and ensuring correct display in the TUI.

## Phase 1: Understand Log Fetching and Raw Content
### Overview
This phase focuses on understanding the exact log content received from the Airflow API and how it's handled *before* any parsing or display logic.

### Changes Required:

#### 1. Temporarily disable log parsing in `airflow/generated_client.go`
**File**: `/home/l/Documents/git/airflow-tui/airflow/generated_client.go`
**Changes**: Comment out the calls to `parseV1LogContent(content)` and the hostname skipping logic within the `GetTaskLog` function.

```go
// In GetTaskLog function:
// content = parseV1LogContent(content)
// lines := strings.Split(content, "\n")
// if len(lines) > 1 && !strings.Contains(lines[0], " ") {
// 	content = strings.Join(lines[1:], "\n")
// }
```

#### 2. Add temporary debug logging in `ui/model.go`
**File**: `/home/l/Documents/git/airflow-tui/ui/model.go`
**Changes**: Add `fmt.Printf` statements in the `logLoadedMsg` handler to print the `msg.content` to `stderr` before it's set to the `logViewport`. This will let us see the processed content in the terminal.

```go
// In Update method, inside case logLoadedMsg:
// fmt.Fprintf(os.Stderr, "DEBUG: Log content received: %s\n", msg.content)
// fmt.Fprintf(os.Stderr, "DEBUG: Log content length: %d\n", len(msg.content))
```

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles successfully: `just build`
- [ ] No new linting errors are introduced: `go vet ./...` (or equivalent `golangci-lint run` if available in `Justfile`)

#### Manual Verification:
- [ ] Run the TUI (`just run`).
- [ ] Navigate to a task's logs.
- [ ] Observe the raw log output in the terminal (stderr).
- [ ] Verify if the "incorrectly shown" behavior changes or becomes clearer with raw output.
- [ ] Take note of the exact format of the raw log content.

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful and the raw log content has been observed, before proceeding to the next phase.

---

## Phase 2: Debug and Correct Log Parsing
### Overview
Based on the raw log content observed in Phase 1, this phase will focus on debugging and correcting the `parseV1LogContent` function and the hostname skipping logic in `airflow/generated_client.go`.

### Changes Required:

#### 1. Re-enable and debug `parseV1LogContent`
**File**: `/home/l/Documents/git/airflow-tui/airflow/generated_client.go`
**Changes**: 
- Re-enable the `content = parseV1LogContent(content)` line.
- Modify `parseV1LogContent` or add debug statements to it to understand why it might be incorrectly parsing the logs. This might involve inspecting the `v1LogRegex` and the `strings.ReplaceAll` logic.

#### 2. Debug hostname skipping logic
**File**: `/home/l/Documents/git/airflow-tui/airflow/generated_client.go`
**Changes**: 
- Re-enable the hostname skipping logic:
```go
// lines := strings.Split(content, "\n")
// if len(lines) > 1 && !strings.Contains(lines[0], " ") {
// 	content = strings.Join(lines[1:], "\n")
// }
```
- Add debug prints to verify if the first line is indeed a hostname and if it's being skipped correctly.

#### 3. Refine `ui/model.go` log display
**File**: `/home/l/Documents/git/airflow-tui/ui/model.go`
**Changes**: Based on the corrected parsing, verify that `m.logViewport.SetContent(m.logContent)` is correctly displaying the content. Adjust `logViewport` dimensions or styling if necessary.

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles successfully: `just build`
- [ ] No new linting errors are introduced: `go vet ./...`

#### Manual Verification:
- [ ] Run the TUI (`just run`).
- [ ] Navigate to a task's logs.
- [ ] Verify that the task logs are displayed correctly, without any parsing errors or extraneous information.
- [ ] Confirm that the hostname line is correctly skipped (if it exists in the raw logs).
- [ ] Ensure that long log lines are handled gracefully by `logViewport`.

---

## Testing Strategy

### Unit Tests:
- No existing Go unit tests were identified. However, we should consider manually creating a test case for `parseV1LogContent` function by creating a temporary `_test.go` file within the `airflow` package to verify its behavior with various log formats. This will ensure that the parsing logic is robust.

### Integration Tests:
- We will rely on manual testing in the TUI, as well as the existing `e2e/test_tui.py` Python e2e tests. We will ensure that the e2e tests related to log viewing (if any exist) pass after the fix. If not, we will identify and update the relevant e2e tests.

### Manual Testing Steps:
1. Start the TUI application: `just run`
2. Select a DAG.
3. Select a DAG run.
4. Select a task instance.
5. Navigate to the "Logs" tab.
6. Observe the log output for correctness and formatting.
7. Test logs for tasks with different "try numbers" (1, 2, 3, etc.) to ensure multi-try logs are handled correctly.
8. Test logs for very long lines or large content to ensure `logViewport` handles them gracefully.

## Performance Considerations
The current approach involves fetching the entire log content and then parsing it. For very large log files, this could be a performance bottleneck. The `token` parameter in `GetTaskLog` suggests pagination, which is already implemented. We should ensure the `logFollow` and `logToken` logic in `ui/model.go` correctly handles continuous loading of paginated logs to avoid fetching the entire log history repeatedly.

## References
- `ui/model.go`
- `airflow/client.go`
- `airflow/generated_client.go`