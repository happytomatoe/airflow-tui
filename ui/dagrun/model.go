package dagrun

import (
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/charmbracelet/bubbles/table"
	uiTheme "github.com/airflow-tui/airflow-tui/ui/theme"
)

// NewDAGRunScreen initializes a new DAGRunScreenModel with preconfigured table
func NewDAGRunScreen() *DAGRunScreenModel {
	t := table.New(
		table.WithColumns(Columns),
		table.WithStyles(table.Styles{
			Header: uiTheme.GetTheme("").TitleStyle,
			Selected: uiTheme.GetTheme("").SelectedStyle,
		}),
	)
	return &DAGRunScreenModel{
		runsTable: t,
	}
}

// SetDAGRuns updates the DAG runs data and refreshes the table rows
func (m *DAGRunScreenModel) SetDAGRuns(runs []airflow.DAGRun) {
	m.runs = runs
	rows := make([]table.Row, 0, len(runs))
	for _, run := range runs {
		runID := derefString(run.DagRunId)
		state := derefString((*string)(run.State))
		executionDate := formatTime(run.ExecutionDate)
		runType := derefString((*string)(run.RunType))
		rows = append(rows, table.Row{
			runID,
			state,
			executionDate,
			runType,
		})
	}
	m.runsTable.SetRows(rows)
}

// derefString converts a *string to string
func derefString(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

// formatTime formats a *time.Time to string
func formatTime(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

// SetClient sets the Airflow client for API calls
func (m *DAGRunScreenModel) SetClient(client airflow.Client) {
	m.client = client
}

// SetSize updates the table dimensions based on terminal size
func (m *DAGRunScreenModel) SetSize(width, height int) {
	m.runsTable.SetWidth(width - 8)
	m.runsTable.SetHeight(height - 10)
}
