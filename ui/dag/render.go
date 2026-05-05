package dag

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/ui/theme"
)

// Render generates the full DAG screen view using original dagListView logic
func (m *DAGScreenModel) Render() string {
	// Original dagListView returns searchView + dagTable.View()
	// This matches the original logic while keeping search handling in main model
	return m.dagTable.View()
}

// scheduleText returns display text for a DAG's schedule (from original model.go:945)
func scheduleText(dag airflow.DAG) string {
	if dag.TimetableDescription != nil && *dag.TimetableDescription != "" {
		return *dag.TimetableDescription
	}
	return "-"
}

// makeRows converts DAG list to table rows (from original model.go:854)
func makeRows(dags []airflow.DAG) []table.Row {
	rows := make([]table.Row, 0, len(dags))
	for _, dag := range dags {
		rows = append(rows, table.Row{
			derefString(dag.DagId),
			scheduleText(dag),
			boolText(dag.IsPaused),
		})
	}
	return rows
}

// Helper functions from original model.go
func derefString(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

func boolText(b *bool) string {
	if b == nil {
		return "-"
	}
	if *b {
		return theme.PausedStyle.Render("yes")
	}
	return theme.ActiveStyle.Render("no")
}
