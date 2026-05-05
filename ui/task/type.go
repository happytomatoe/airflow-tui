package task

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/airflow-tui/airflow-tui/airflow"
)

type TaskScreenModel struct {
	taskTable     table.Model
	taskInstances []airflow.TaskInstance
	client        airflow.Client
	dagID         string
	dagRunID      string
	showGantt     bool
}
