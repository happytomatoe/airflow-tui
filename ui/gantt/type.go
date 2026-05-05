package gantt

import (
	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

type GanttDialogModel struct {
	ganttChart GanttChart
	visible     bool
}

func NewGanttDialog() *GanttDialogModel {
	return &GanttDialogModel{
		ganttChart: NewGanttChart(),
	}
}

func (m *GanttDialogModel) Show() {
	m.visible = true
}

func (m *GanttDialogModel) Hide() {
	m.visible = false
}

func (m *GanttDialogModel) IsVisible() bool {
	return m.visible
}

func (m *GanttDialogModel) SetTasks(tasks []airflow.TaskInstance) {
	m.ganttChart.SetTasks(tasks)
}

func (m *GanttDialogModel) SetWidth(width int) {
	m.ganttChart.SetWidth(width)
}

// View implements tea.Model interface
func (m *GanttDialogModel) View() string {
	if !m.visible {
		return ""
	}
	return m.ganttChart.View()
}

var ganttHeaderStyle = lipgloss.NewStyle().
	Foreground(theme.GetTheme("").AccentColor).
	Bold(true).
	Padding(0, 0, 1, 0)
