package task

import (
	"github.com/airflow-tui/airflow-tui/ui/theme"
	"github.com/airflow-tui/airflow-tui/ui/gantt"
)

func (m *TaskScreenModel) Render() string {
	if len(m.taskInstances) == 0 {
		return theme.GetTheme("").MutedStyle.Render("No task instances found")
	}
	if m.showGantt {
		g := gantt.NewGanttChart()
		g.SetTasks(m.taskInstances)
		return g.View()
	}
	return m.taskTable.View()
}
