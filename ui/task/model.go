package task

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/airflow-tui/airflow-tui/airflow"
)

func NewTaskScreen() *TaskScreenModel {
	taskTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "Task ID", Width: 40},
			{Title: "State", Width: 12},
			{Title: "Try #", Width: 6},
			{Title: "Duration", Width: 10},
		}),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	taskStyles := table.DefaultStyles()
	taskStyles.Header = taskStyles.Header.Bold(true)
	taskStyles.Selected = taskStyles.Selected.Bold(true)
	taskTable.SetStyles(taskStyles)

	return &TaskScreenModel{
		taskTable: taskTable,
	}
}

func (m *TaskScreenModel) SetTasks(tasks []airflow.TaskInstance) {
	m.taskInstances = tasks
	rows := make([]table.Row, 0, len(tasks))
	for _, task := range tasks {
		state := ""
		if task.State != nil {
			state = string(*task.State)
		}
		rows = append(rows, table.Row{
			derefString(task.TaskId),
			state,
			formatTryNumber(task),
			formatDuration(task),
		})
	}
	m.taskTable.SetRows(rows)
}

func (m *TaskScreenModel) SetClient(client airflow.Client) {
	m.client = client
}

func (m *TaskScreenModel) SetDagInfo(dagID, dagRunID string) {
	m.dagID = dagID
	m.dagRunID = dagRunID
}

func (m *TaskScreenModel) ToggleGantt() {
	m.showGantt = !m.showGantt
}

func (m *TaskScreenModel) SetSize(width, height int) {
	m.taskTable.SetWidth(width - 8)
	m.taskTable.SetHeight(height - 10)
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatTryNumber(task airflow.TaskInstance) string {
	return "1" // TODO: extract from task
}

func formatDuration(task airflow.TaskInstance) string {
	return "" // TODO: calculate from start/end time
}
