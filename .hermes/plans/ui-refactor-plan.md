# Airflow-TUI UI Refactoring Plan
## Goal
Refactor the airflow-tui UI codebase to follow the per-screen file organization pattern demonstrated in [superfile](https://github.com/yorukot/superfile), where each screen has its own directory with separated concerns (type, model, render, update). Standardize on "screen" terminology (not "panel").

## Domain Language (see CONTEXT.md for full definitions)
- **Screen**: Tab-based main content (dag, dagrun, task, log, config)
- **Dialog**: Overlay component (confirm, help, search, gantt)
- All async message types kept in shared location (main model.go or ui/messages.go)
- Screen models stored as **pointers** in main Model
- Dialog visibility flags live in dialog models (not main Model)

## Current State Analysis

### Existing UI Structure (in \`ui/\`)
| File | Lines | Purpose |
|------|-------|---------|
| \`model.go\` | 1444 | Main \`Model\` struct, all async message types, core TUI logic |
| \`gantt.go\` | 186 | \`GanttChart\` struct and rendering |
| \`confirm.go\` | 44 | \`renderConfirm()\` function, confirm dialog styles |
| \`help.go\` | 120 | \`HelpItem\` struct, \`renderHelp()\`, per-tab help |
| \`tabs.go\` | 131 | \`Tab\` enum, \`TabBar\` struct, tab switching logic |
| \`search_dialog.go\` | 32 | \`searchDialogModel\` struct, \`renderSearchDialog()\` |
| \`theme.go\` | 28 | Base style definitions |
| \`themes.go\` | 208 | \`Theme\` struct, 6 theme presets |

### Current Screens and Dialogs
**Screens (tab-based):**
1. **DAG Screen** - Lists DAGs in a table
2. **DAG Run Screen** - Lists DAG runs for a selected DAG
3. **Task Screen** - Lists task instances for a DAG run
4. **Log Screen** - Displays logs for a task instance
5. **Config Screen** - Displays configuration in a table

**Dialogs (overlays):**
6. **Confirm Dialog** - Modal dialog for confirming actions
7. **Help Dialog** - Overlay showing keyboard shortcuts
8. **Search Dialog** - Dialog for filtering/searching
9. **Gantt Dialog** - Overlay showing Gantt chart for task instances

### Pain Points
- \`model.go\` is a 1444-line monolith
- All async message types defined in \`model.go\`
- Hard to navigate and maintain
- No clear separation between screen-specific and shared code

## Proposed Refactoring Steps

### Step 1: Create Directory Structure
Create per-screen directories under \`ui/\` using component names (matching superfile pattern):

\`\`\`
ui/
├── dag/               # DAG list screen
├── dagrun/            # DAG run list screen
├── task/              # Task instance screen
├── log/               # Log viewer screen
├── config/            # Config table screen
├── confirm/           # Confirm dialog
├── help/              # Help dialog
├── search/            # Search dialog
├── gantt/             # Gantt chart dialog
├── tabs/              # Tab bar shared component
├── theme/             # Shared styles and themes
├── nav/               # Navigation struct (avoids import cycles)
└── model.go           # Main coordinator model (slimmed down to ~300 lines)
\`\`\`

**Package naming**: Match directory names (standard Go convention). E.g., \`ui/dag/\` → package \`dag\`.

### Step 2: Create Navigation Package
**Directory: \`ui/nav/\`**
- Move \`navigation\` struct here from \`model.go\`
- Export fields: \`Server\`, \`Dag\`, \`DagRun\`, \`Task\`, \`LogType\`
- Both main \`ui/\` and screen packages import \`ui/nav/\` to avoid import cycles

### Step 3: Extract Theme/Styles
**Directory: \`ui/theme/\`**
- Move \`theme.go\` and \`themes.go\` here
- Export \`GetStyles()\` function that returns current theme's styles
- All screens import \`ui/theme\` to access shared styles (not local styles)

### Step 4: Extract Tab Component
**Directory: \`ui/tabs/\`**
- Move \`Tab\` enum, \`TabBar\` struct, \`TabKeyMap\` from \`tabs.go\`
- Separate into: \`type.go\`, \`model.go\` (constructor, state), \`render.go\`, \`update.go\`
- Main Model embeds \`tabBar *tabs.TabBar\` (pointer)

### Step 5: Extract DAG Screen
**Directory: \`ui/dag/\`**
- \`type.go\`: \`DAGScreenModel\` struct (contains \`dagTable table.Model\`, \`dags []airflow.DAG\`)
- \`model.go\`: \`NewDAGScreen()\`, \`SetDAGs()\`, \`SetClient(client)\`, \`SetConfig(cfg)\`
- \`render.go\`: \`Render() string\` method (convert from \`viewDAGPanel()\`)
- \`update.go\`: \`FetchDAGs() tea.Cmd\` async function in \`fetch.go\`
- \`consts.go\`: Column definitions for DAG table

### Step 6: Extract DAG Run Screen
**Directory: \`ui/dagrun/\`**
- \`type.go\`: \`DAGRunScreenModel\` struct
- \`model.go\`: \`NewDAGRunScreen()\`, \`SetDAGRuns()\`, \`SetClient(client)\`
- \`render.go\`: \`Render() string\` method
- \`fetch.go\`: \`FetchDAGRuns(client, dagID) tea.Cmd\`
- \`consts.go\`: Column definitions for DAG run table

### Step 7: Extract Task Screen
**Directory: \`ui/task/\`**
- \`type.go\`: \`TaskScreenModel\` struct (contains \`taskTable table.Model\`, \`taskInstances []airflow.TaskInstance\`)
- \`model.go\`: \`NewTaskScreen()\`, \`SetTasks()\`, \`SetClient(client)\`
- \`render.go\`: \`Render() string\` method (task table only)
- \`fetch.go\`: \`FetchTaskInstances(client, dagID, dagRunID) tea.Cmd\`

### Step 8: Extract Gantt Dialog
**Directory: \`ui/gantt/\`**
- \`type.go\`: \`GanttDialogModel\` struct (contains \`ganttChart GanttChart\`, \`visible bool\`)
- \`model.go\`: \`NewGanttDialog()\`, \`SetTasks()\`, \`Show()\`, \`Hide()\`, \`Visible() bool\`, \`SetClient(client)\`
- \`render.go\`: \`Render() string\` method
- \`gantt.go\`: \`GanttChart\` struct and methods (moved from \`gantt.go\`)

### Step 9: Extract Log Screen
**Directory: \`ui/log/\`**
- \`type.go\`: \`LogScreenModel\` struct
- \`model.go\`: \`NewLogScreen()\`, \`SetLogs()\`, \`SetClient(client)\`, log follow logic
- \`render.go\`: \`Render() string\` method using \`logViewport\`
- \`fetch.go\`: \`FetchLogs(client, dagID, dagRunID, taskID) tea.Cmd\`

### Step 10: Extract Config Screen
**Directory: \`ui/config/\`**
- \`type.go\`: \`ConfigScreenModel\` struct
- \`model.go\`: \`NewConfigScreen()\`, \`SetConfig()\`
- \`render.go\`: \`Render() string\` method

### Step 11: Extract Dialog Components

**Directory: \`ui/confirm/\`**
- \`type.go\`: \`ConfirmDialogModel\` struct (contains \`title\`, \`message\`, \`callback func() tea.Cmd\`, \`visible bool\`)
- \`model.go\`: \`NewConfirmDialog()\`, \`ShowConfirm(title, msg, callback)\`, \`Hide()\`, \`Visible() bool\`, \`Confirm() tea.Cmd\`
- \`render.go\`: \`Render() string\` method (convert from \`renderConfirm()\`)

**Directory: \`ui/help/\`**
- \`type.go\`: \`HelpDialogModel\` struct (contains \`visible bool\`)
- \`model.go\`: \`NewHelpDialog()\`, \`Show()\`, \`Hide()\`, \`Visible() bool\`
- \`render.go\`: \`Render() string\` method (convert from \`renderHelp()\`)

**Directory: \`ui/search/\`**
- \`type.go\`: \`SearchDialogModel\` struct (contains \`visible bool\`, \`input textinput.Model\`)
- \`model.go\`: \`NewSearchDialog()\`, \`Show()\`, \`Hide()\`, \`Visible() bool\`
- \`render.go\`: \`Render() string\` method (convert from \`renderSearchDialog()\`)

### Step 12: Create NewModel() Constructor
**In \`ui/model.go\`:**
\`\`\`go
func NewModel(cfg config.Config) *Model {
    m := &Model{
        cfg:        cfg,
        activeName: activeClient(cfg),
        loading:    true,
        // Initialize screen pointers
        dagScreen:      dag.NewDAGScreen(),
        dagRunScreen:   dagrun.NewDAGRunScreen(),
        taskScreen:     task.NewTaskScreen(),
        logScreen:      log.NewLogScreen(),
        configScreen:   config.NewConfigScreen(),
        // Initialize dialogs
        confirmDialog:  confirm.NewConfirmDialog(),
        helpDialog:     help.NewHelpDialog(),
        searchDialog:   search.NewSearchDialog(),
        ganttDialog:    gantt.NewGanttDialog(),
        // Shared components
        tabBar:         tabs.NewTabBar(),
        nav:            nav.Navigation{},
    }
    
    // Set client and config on screens
    client, _ := activeClientWithClient(cfg)
    m.dagScreen.SetClient(client)
    m.dagRunScreen.SetClient(client)
    // ... etc
    
    return m
}
\`\`\`

### Step 13: Refactor Main Model
**Slimmed \`ui/model.go\` (~300 lines):**
\`\`\`go
type Model struct {
    cfg        config.Config
    client     airflow.Client  
    activeName string
    loading    bool
    err        error
    connected  bool
    
    // Screen pointers
    dagScreen      *dag.DAGScreenModel
    dagRunScreen   *dagrun.DAGRunScreenModel
    taskScreen     *task.TaskScreenModel
    logScreen      *log.LogScreenModel
    configScreen   *config.ConfigScreenModel
    
    // Dialog pointers
    confirmDialog  *confirm.ConfirmDialogModel
    helpDialog     *help.HelpDialogModel
    searchDialog   *search.SearchDialogModel
    ganttDialog    *gantt.GanttDialogModel
    
    // Shared
    tabBar        *tabs.TabBar
    nav           nav.Navigation
    
    activeScreen  int  // renamed from panel to screen
}
\`\`\`

**Key removals from main Model:**
- \`width\`, \`height\` fields (screens handle WindowSizeMsg internally)
- \`spinner spinner.Model\` (loading indicated by text only)
- \`showHelp bool\`, \`searching bool\` (now in dialog models)

### Step 14: Implement Main Update() Method
**Main Update() handles ALL keys and async messages (Option A):**
\`\`\`go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // 1. Global keys first (quit, tab switching)
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "tab":
            m.switchScreen()
            return m, nil
        // ... handle help toggle, search toggle
        }
    case tea.WindowSizeMsg:
        // Route to active screen
        m.dagScreen.HandleWindowSize(msg) // or similar
        return m, nil
    }
    
    // 2. Async messages - main Update handles these
    switch msg := msg.(type) {
    case dagsLoadedMsg:
        m.dagScreen.SetDAGs(msg.dags)
        m.loading = false
        return m, nil
    case dagRunsLoadedMsg:
        m.dagRunScreen.SetDAGRuns(msg.runs)
        m.loading = false
        return m, nil
    // ... other async messages
    }
    
    return m, nil
}
\`\`\`

### Step 15: Implement Main View() Method
**When dialog visible, ONLY render the dialog:**
\`\`\`go
func (m Model) View() string {
    // Check dialogs first (dialogs replace screen content when visible)
    if m.confirmDialog != nil && m.confirmDialog.Visible() {
        return m.confirmDialog.Render()
    }
    if m.helpDialog != nil && m.helpDialog.Visible() {
        return m.helpDialog.Render()
    }
    // ... other dialogs
    
    // Render active screen
    switch m.activeScreen {
    case dagScreen:
        return m.dagScreen.Render()
    case dagRunScreen:
        return m.dagRunScreen.Render()
    // ... other screens
    }
}
\`\`\`

### Step 16: Update Import Paths
- Change all imports from \`github.com/airflow-tui/airflow-tui/ui\` to new packages
- Example: \`import "github.com/airflow-tui/airflow-tui/ui/dag"\`
- Update \`main.go\` to use \`ui.NewModel(cfg)\` instead of creating Model directly

## File References
- Main model to refactor: \`ui/model.go\` (1444 lines)
- Example superfile component: \`/tmp/superfile/src/internal/ui/helpmenu/\`
- ADR: \`docs/adr/0001-per-screen-refactoring.md\`
- Context: \`CONTEXT.md\`

## Rollout Strategy
**All-at-once on feature branch** (not incremental). Complete all extractions in one branch, then slim down \`model.go\`. This results in a larger PR but avoids intermediate states.

## Automated Verification
1. **Build check**: \`go build ./...\` passes
2. **Linting**: \`go fmt ./...\` and \`go vet ./...\` pass
3. **E2E tests**: \`just e2e-test\` passes (requires running Airflow instance)
4. **Unit tests**: Note: No Go unit tests currently exist per AGENTS.md

## Manual Verification  
1. **UI functionality**: All screens render correctly (DAG list, DAG runs, tasks, logs, config)
2. **Navigation**: Tab switching, screen navigation works
3. **Dialogs**: Confirm, Help, Search, Gantt dialogs function properly
4. **Themes**: All 6 themes apply correctly
5. **Log follow**: Log auto-follow mode works in log screen
6. **No behavior change**: UI behaves identically after refactoring

## Out-of-Scope
- **Adding new features**: Pure refactoring effort
- **Changing behavior**: UI should behave identically
- **Adding unit tests**: None exist currently
- **Refactoring \`airflow/\` package**: Only UI code is in scope
- **Changing e2e tests**: Minor import path updates only if needed

## Risks and Mitigations
| Risk | Mitigation |
|------|------------|
| Breaking UI behavior | All-at-once refactoring on feature branch, thorough manual verification |
| Import cycle issues | Navigation struct in separate \`ui/nav/\` package |
| E2E test failures | Run \`just e2e-test\` after complete refactoring |
| Styles not rendering | Verify theme imports after moving to \`ui/theme/\` |

## Success Criteria
1. Each screen/dialog has its own directory under \`ui/\` with separated concerns
2. Main \`model.go\` reduced from 1444 lines to <300 lines
3. All automated and manual verification steps pass
4. No change in UI behavior from user perspective
5. Dialogs render as overlays (replacing screen content when visible)
