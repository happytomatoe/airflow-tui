     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|     1|# CONTEXT.md
     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|     2|
     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|     3|## Domain Language
     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|     4|
     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|     5|### Screen
     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|     6|A distinct visual area in the TUI that is the main content for a tab. Screens replace each other when switching tabs, controlled by the active screen enum in the main model.
     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|     7|
     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|     8|Screens in airflow-tui:
     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|     9|- **DAG Screen** - Lists DAGs in a table (tab: TabDags)
    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|    10|- **DAG Run Screen** - Lists DAG runs for a selected DAG (tab: TabRuns)
    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|    11|- **Task Screen** - Lists task instances for a DAG run (tab: TabTasks)
    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|    12|- **Log Screen** - Displays logs for a task instance (tab: TabLogs)
    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|    13|- **Config Screen** - Displays configuration in a table (tab: TabConfig)
    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|    14|
    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|    15|### Dialog
    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|    16|An overlay component that renders on top of the active screen, controlled by a visibility flag (e.g., `showHelp`, `confirm != nil`). Dialogs do not replace the underlying screen.
    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|    17|
    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|    18|Dialogs in airflow-tui:
    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|    19|- **Confirm Dialog** - Modal dialog for confirming actions, has a callback on confirm
    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|    20|- **Help Dialog** - Overlay showing keyboard shortcuts per tab
    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|    21|- **Search Dialog** - Dialog for filtering/searching the current screen's content
    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|    22|
    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|    23|### Async Messages
    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|    24|All async message types (e.g., `dagsLoadedMsg`, `dagRunsLoadedMsg`) are kept in a shared location (either `ui/model.go` or `ui/messages.go`) because the main `Model.Update()` uses a type switch that needs to reference all message types directly. Splitting them per screen adds import complexity without benefit.
    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|    25|
    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|    26|### Naming Note
    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|    27|The codebase currently uses `type panel int` and `dagPanel`, `taskPanel` etc. in `ui/model.go`. During refactoring:
    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|    28|- Enum will be renamed to `type screen int` with `dagScreen`, `taskScreen` etc.
    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|    29|- Directory names will reflect component purpose (see next question)
    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|    30|
    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|    31|
    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|    32|### Screen Model Pointers
    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|    33|Screen models are stored as **pointers** in the main Model struct (e.g., `dagScreen *dag.DAGScreenModel`). This allows sub-Update() methods to modify screen state in-place without returning a modified copy. Shared components like Tab Bar also use pointers.
    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|    34|
    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|    35|
    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|    36|### Async Message Routing
    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|    37|Each screen's Update() method handles its own async messages (e.g., `dagRunsLoadedMsg` handled by `dagRunScreen.Update()`). The main Model.Update() routes messages to the responsible screen based on message type.
    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|    38|
    39|    39|    39|    39|    39|    39|    39|    39|    39|    39|    39|    39|    39|    39|    39|    39|
    40|    40|    40|    40|    40|    40|    40|    40|    40|    40|    40|    40|    40|    40|    40|    40|### Shared Styles
    41|    41|    41|    41|    41|    41|    41|    41|    41|    41|    41|    41|    41|    41|    41|    41|All screens import `ui/theme` package to access shared styles (e.g., `theme.GetStyles()`). The theme package exports base styles and theme presets. Each screen does NOT define its own styles locally — consistent theming is a global concern.
    42|    42|    42|    42|    42|    42|    42|    42|    42|    42|    42|    42|    42|    42|    42|    42|
    43|    43|    43|    43|    43|    43|    43|    43|    43|    43|    43|    43|    43|    43|    43|
    44|    44|    44|    44|    44|    44|    44|    44|    44|    44|    44|    44|    44|    44|    44|### Client and Config Access
    45|    45|    45|    45|    45|    45|    45|    45|    45|    45|    45|    45|    45|    45|    45|Screens receive the Airflow client and config at initialization via `SetClient(client)` and `SetConfig(cfg)` methods. The client is stored as a pointer field on each screen model. This allows screens to trigger their own async API operations (e.g., `dagScreen.FetchDAGs()`) without the main Model intermediating every API call.
    46|    46|    46|    46|    46|    46|    46|    46|    46|    46|    46|    46|    46|    46|    46|
    47|    47|    47|    47|    47|    47|    47|    47|    47|    47|    47|    47|    47|    47|
    48|    48|    48|    48|    48|    48|    48|    48|    48|    48|    48|    48|    48|    48|### Confirm Dialog Callback
    49|    49|    49|    49|    49|    49|    49|    49|    49|    49|    49|    49|    49|    49|The callback for the confirm dialog is passed at show time via `ShowConfirm(title, msg, callback)`. This allows the caller (main Update or another screen) to define what happens on confirm each time the dialog is shown, without the dialog needing to know about specific actions. The callback returns a `tea.Cmd` that the main Update() executes.
    50|    50|    50|    50|    50|    50|    50|    50|    50|    50|    50|    50|    50|    50|
    51|    51|    51|    51|    51|    51|    51|    51|    51|    51|    51|    51|    51|
    52|    52|    52|    52|    52|    52|    52|    52|    52|    52|    52|    52|    52|### Refactoring Rollout Strategy
    53|    53|    53|    53|    53|    53|    53|    53|    53|    53|    53|    53|    53|Use a feature branch and do the entire refactoring at once (rather than incremental extraction). This is faster but results in a larger, harder-to-review PR. All screens, dialogs, and shared components will be extracted in one branch, then the main `model.go` will be slimmed down to a coordinator.
    54|    54|    54|    54|    54|    54|    54|    54|    54|    54|    54|    54|    54|
    55|    55|    55|    55|    55|    55|    55|    55|    55|    55|    55|    55|
    56|    56|    56|    56|    56|    56|    56|    56|    56|    56|    56|    56|### Global State (Loading, Error, Connection)
    57|    57|    57|    57|    57|    57|    57|    57|    57|    57|    57|    57|Loading (`loading bool`), error (`err error`), connection status (`connected bool`), and spinner (`spinner spinner.Model`) live in the main Model as global state. Screens do NOT have their own loading/error fields. Instead, screens report errors via return values or messages to the main Update(), which then sets the global state.
    58|    58|    58|    58|    58|    58|    58|    58|    58|    58|    58|    58|
    59|    59|    59|    59|    59|    59|    59|    59|    59|    59|    59|
    60|    60|    60|    60|    60|    60|    60|    60|    60|    60|    60|### Dialog Rendering
    61|    61|    61|    61|    61|    61|    61|    61|    61|    61|    61|When a dialog is visible, the main `View()` ONLY renders the dialog (not the underlying screen content). This is the simpler approach for terminal TUIs - no complex overlay rendering needed. The dialog completely replaces the screen content while visible.
    62|    62|    62|    62|    62|    62|    62|    62|    62|    62|    62|
    63|    63|    63|    63|    63|    63|    63|    63|    63|    63|
    64|    64|    64|    64|    64|    64|    64|    64|    64|    64|### Navigation Struct
    65|    65|    65|    65|    65|    65|    65|    65|    65|    65|The `navigation` struct (tracking server, dag, dagRun, task, logType) lives in the main Model as global navigation state. It represents cross-screen state - selecting a DAG in DAG Screen affects what appears in DAG Run Screen. Screens can READ navigation via their methods, but the main Model OWNS it.
    66|    66|    66|    66|    66|    66|    66|    66|    66|    66|
    67|    67|    67|    67|    67|    67|    67|    67|    67|
    68|    68|    68|    68|    68|    68|    68|    68|    68|### Model Initialization
    69|    69|    69|    69|    69|    69|    69|    69|    69|A `NewModel()` constructor initializes all screen pointers, sets client/config on each screen, and returns the fully initialized main Model. The Bubble Tea `Init()` method then just returns the initial command (e.g., `fetchDAGs()`). This avoids the issue of `Init()` receiving the Model by value.
    70|    70|    70|    70|    70|    70|    70|    70|    70|
    71|    71|    71|    71|    71|    71|    71|    71|
    72|    72|    72|    72|    72|    72|    72|    72|### Package Naming
    73|    73|    73|    73|    73|    73|    73|    73|Go package names match their directory names (standard Go convention). Examples:
    74|    74|    74|    74|    74|    74|    74|    74|- `ui/dag/` → package `dag`
    75|    75|    75|    75|    75|    75|    75|    75|- `ui/task/` → package `task`
    76|    76|    76|    76|    76|    76|    76|    76|- `ui/confirm/` → package `confirm`
    77|    77|    77|    77|    77|    77|    77|    77|- `ui/tabs/` → package `tabs`
    78|    78|    78|    78|    78|    78|    78|    78|- `ui/theme/` → package `theme`
    79|    79|    79|    79|    79|    79|    79|    79|
    80|    80|    80|    80|    80|    80|    80|    80|Imports look like: `import "github.com/airflow-tui/airflow-tui/ui/dag"`
    81|    81|    81|    81|    81|    81|    81|    81|
    82|    82|    82|    82|    82|    82|    82|
    83|    83|    83|    83|    83|    83|    83|### Render Functions as Methods
    84|    84|    84|    84|    84|    84|    84|All standalone `render*()` functions (e.g., `renderConfirm()`, `renderHelp()`) become `Render() string` methods on their respective screen/dialog structs. This is cleaner OOP and matches the superfile pattern. The functions currently take `*Model` (main model) but after refactoring, each struct holds its own state.
    85|    85|    85|    85|    85|    85|    85|
    86|    86|    86|    86|    86|    86|
    87|    87|    87|    87|    87|    87|### Dialog Visibility Flags
    88|    88|    88|    88|    88|    88|Each dialog owns its visibility state via a `visible bool` field (or `Visible()` method). The main Model does NOT have separate `showHelp bool`, `searching bool` fields. Instead, it checks `m.helpDialog.Visible()`, `m.confirmDialog.Visible()`, etc. This reduces main Model fields and makes each dialog self-contained.
    89|    89|    89|    89|    89|    89|
    90|    90|    90|    90|    90|
    91|    91|    91|    91|    91|### Async Command Functions
    92|    92|    92|    92|    92|Async command functions (e.g., `fetchDAGs()`, `fetchDAGRuns()`) live in their respective screen packages as `fetch.go` files. Each screen owns its async operations. The main Update() calls these functions when handling user actions (e.g., selecting a DAG calls `dagrun.FetchDAGRuns()`).
    93|    93|    93|    93|    93|
    94|    94|    94|
    95|    95|    95|### Screen Dimensions (Width/Height)
    96|    96|    96|The main Model does NOT store `width` and `height` fields. Instead, `tea.WindowSizeMsg` is routed to the active screen's `Update()`, and each screen stores its own dimensions if needed. In the current codebase, `m.width` was only used once (for Gantt chart width), so this is safe to remove.
    97|    97|    97|
    98|    98|
    99|    99|### Key Handling
   100|   100|The main `Model.Update()` handles ALL key presses (including screen-specific keys). Screens do NOT have `HandleKey()` methods. Instead, main Update() has explicit key handling per screen (switch on `m.activeScreen` inside key handler). Screens expose methods like `ToggleGantt()`, `SelectTask()`, etc., which main Update() calls.
   101|   101|