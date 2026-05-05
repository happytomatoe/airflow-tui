package dagrun

import (
	uiTheme "github.com/airflow-tui/airflow-tui/ui/theme"
)

// Render returns the string representation of the DAG Run screen
func (m *DAGRunScreenModel) Render() string {
	if len(m.runs) == 0 {
		return uiTheme.GetTheme("").MutedStyle.Render("No DAG runs found")
	}
	return m.runsTable.View()
}
