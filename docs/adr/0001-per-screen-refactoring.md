# Per-screen directory refactoring

We refactored the monolithic `ui/model.go` (1444 lines) into per-screen directories (`ui/dag/`, `ui/task/`, etc.) following the pattern from [superfile](https://github.com/yorukot/superfile). Each screen gets its own directory with separated concerns: `type.go`, `model.go`, `render.go`, `update.go`. This makes the codebase easier to navigate and maintain as new screens are added.

## Status
accepted

## Considered Options
- **Keep monolith**: Keep all logic in `ui/model.go`. Rejected because it's hard to navigate and maintain as the app grows.
- **Per-screen directories (chosen)**: Follow superfile's pattern. Provides clear separation of concerns and matches established conventions in similar Bubble Tea apps.
- **Split by feature but different structure**: Rejected because superfile's pattern is proven and familiar.

## Consequences
- All imports change: screens now live under `ui/<screen>/` packages.
- Main `Model` becomes a coordinator (~300 lines) that embeds screen pointers and handles global concerns (loading, errors, navigation).
- Dialogs (confirm, help, search, gantt) follow the same per-directory pattern as screens.
- Async message handling remains in main `Update()` (Option A) for simplicity.
- Key handling remains in main `Update()` for simplicity (Option A).
- `spinner.Model` removed; loading indicated by text only.
- `width`/`height` fields removed from main Model; screens handle `tea.WindowSizeMsg` internally if needed.
