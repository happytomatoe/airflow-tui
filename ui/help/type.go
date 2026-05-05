package help

import (
	"fmt"

	"github.com/airflow-tui/airflow-tui/ui/tabs"
	"github.com/airflow-tui/airflow-tui/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

type HelpItem struct {
	Key   string
	Desc  string
	Scope string
}

type HelpDialogModel struct {
	visible bool
}

func NewHelpDialog() *HelpDialogModel {
	return &HelpDialogModel{}
}

func (m *HelpDialogModel) Show() {
	m.visible = true
}

func (m *HelpDialogModel) Hide() {
	m.visible = false
}

func (m *HelpDialogModel) IsVisible() bool {
	return m.visible
}

func (m *HelpDialogModel) Render(tab tabs.Tab) string {
	if !m.visible {
		return ""
	}

	items := helpItemsForTab(tab)

	var lines []string
	lines = append(lines, theme.GetTheme("").TitleStyle.Render("Help"))

	for _, item := range items {
		key := theme.GetTheme("").TitleStyle.Padding(0, 1, 0, 0).Render(fmt.Sprintf("%-6s", item.Key))
		desc := theme.GetTheme("").MutedStyle.Render(item.Desc)
		lines = append(lines, key+desc)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.GetTheme("").BorderColor).
		Padding(1, 2)

	return boxStyle.Render(content)
}

func helpItemsForTab(tab tabs.Tab) []HelpItem {
	switch tab {
	case tabs.TabConfig:
		return []HelpItem{
			{"q", "quit", "all"},
			{"1-5", "switch tab", "all"},
			{"←/→", "prev/next tab", "all"},
			{"?", "toggle help", "all"},
		}
	case tabs.TabDags:
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
	case tabs.TabRuns:
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
	case tabs.TabTasks:
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
	case tabs.TabLogs:
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
