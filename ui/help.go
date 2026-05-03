package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type HelpItem struct {
	Key   string
	Desc  string
	Scope string
}

var (
	helpTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			Padding(0, 0, 1, 0)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Padding(0, 1, 0, 0)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	helpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("212")).
			Padding(1, 2)
)

func renderHelp(tab Tab) string {
	items := helpItemsForTab(tab)

	var lines []string
	lines = append(lines, helpTitleStyle.Render("Help"))
	for _, item := range items {
		key := helpKeyStyle.Render(fmt.Sprintf("%-6s", item.Key))
		desc := helpDescStyle.Render(item.Desc)
		lines = append(lines, key+desc)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return helpBoxStyle.Render(content)
}

func helpItemsForTab(tab Tab) []HelpItem {
	switch tab {
	case TabConfig:
		return []HelpItem{
			{"q", "quit", "all"},
			{"1-5", "switch tab", "all"},
			{"←/→", "prev/next tab", "all"},
			{"?", "toggle help", "all"},
		}
	case TabDags:
		return []HelpItem{
			{"q", "quit", "all"},
			{"r", "refresh DAGs", "dags"},
			{"/", "filter/search", "dags"},
			{"j/k", "move up/down", "dags"},
			{"enter", "view DAG runs", "dags"},
			{"p", "pause/unpause DAG", "dags"},
			{"t", "trigger DAG run", "dags"},
			{"v", "view DAG code", "dags"},
			{"o", "open in browser", "dags"},
			{"1-5", "switch tab", "all"},
			{"←/→", "prev/next tab", "all"},
			{"?", "toggle help", "all"},
		}
	case TabRuns:
		return []HelpItem{
			{"esc", "back to DAGs", "runs"},
			{"q", "quit", "all"},
			{"r", "refresh runs", "runs"},
			{"j/k", "move up/down", "runs"},
			{"enter", "view tasks", "runs"},
			{"V", "visual mode (multi-select)", "runs"},
			{"c", "clear selected runs", "runs"},
			{"m", "mark runs (success/failed)", "runs"},
			{"d", "show dependency graph", "runs"},
			{"1-5", "switch tab", "all"},
			{"←/→", "prev/next tab", "all"},
			{"?", "toggle help", "all"},
		}
	case TabTasks:
		return []HelpItem{
			{"esc", "back to runs", "tasks"},
			{"q", "quit", "all"},
			{"r", "refresh tasks", "tasks"},
			{"j/k", "move up/down", "tasks"},
			{"enter", "view logs", "tasks"},
			{"V", "visual mode (multi-select)", "tasks"},
			{"c", "clear selected tasks", "tasks"},
			{"m", "mark tasks", "tasks"},
			{"o", "open in Airflow UI", "tasks"},
			{"g", "view gantt chart", "tasks"},
			{"1-5", "switch tab", "all"},
			{"←/→", "prev/next tab", "all"},
			{"?", "toggle help", "all"},
		}
	case TabLogs:
		return []HelpItem{
			{"esc", "back to tasks", "logs"},
			{"q", "quit", "all"},
			{"j/k", "scroll up/down", "logs"},
			{"gg", "jump to top", "logs"},
			{"G", "jump to bottom", "logs"},
			{"F", "toggle follow mode", "logs"},
			{"1-9", "switch try number", "logs"},
			{"o", "open in Airflow UI", "logs"},
			{"1-5", "switch tab", "all"},
			{"←/→", "prev/next tab", "all"},
			{"?", "toggle help", "all"},
		}
	}
	return []HelpItem{}
}
