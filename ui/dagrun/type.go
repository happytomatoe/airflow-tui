package dagrun

import (
	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/charmbracelet/bubbles/table"
)

// DAGRunScreenModel represents the DAG Run screen state
type DAGRunScreenModel struct {
	runsTable table.Model
	runs      []airflow.DAGRun
	client    airflow.Client
}
