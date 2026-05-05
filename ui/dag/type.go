package dag

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/airflow-tui/airflow-tui/airflow"
)

// DAGScreenModel encapsulates the DAG listing screen state
type DAGScreenModel struct {
	dagTable table.Model
	dags     []airflow.DAG
	client   airflow.Client
}
