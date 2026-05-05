package dag

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/ui/theme"
)

// NewDAGScreen initializes a new DAG screen model with preconfigured table
func NewDAGScreen() *DAGScreenModel {
	dagTable := table.New(
		table.WithColumns(DagTableColumns),
		table.WithStyles(table.Styles{
			Header: theme.GetTheme("").TitleStyle,
			Selected: theme.GetTheme("").SelectedStyle,
		}),
	)

	return &DAGScreenModel{
		dagTable: dagTable,
	}
}

// SetDAGs updates the list of DAGs and refreshes the table rows
func (m *DAGScreenModel) SetDAGs(dags []airflow.DAG) {
	m.dags = dags
	rows := makeRows(dags)
	m.dagTable.SetRows(rows)
	if len(rows) > 0 {
		m.dagTable.SetCursor(0)
	}
}

// SetClient sets the Airflow API client for DAG operations
func (m *DAGScreenModel) SetClient(client airflow.Client) {
	m.client = client
}

// SetSize updates the table dimensions based on terminal size
func (m *DAGScreenModel) SetSize(width, height int) {
	m.dagTable.SetWidth(width - 8)
	m.dagTable.SetHeight(height - 10)
}

// SetConfig applies configuration updates to the DAG screen (placeholder)
func (m *DAGScreenModel) SetConfig() {
	// Configuration application logic would go here
}
