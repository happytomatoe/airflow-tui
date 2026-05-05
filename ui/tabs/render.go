package tabs

import (
	"github.com/airflow-tui/airflow-tui/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

func (t *TabBar) View() string {
	var tabs []string
	for i := range tabNames {
		tabStr := tabNames[i]
		if Tab(i) == t.active {
			tabs = append(tabs, theme.GetTheme("").TabActiveStyle.Render(tabStr))
		} else {
			tabs = append(tabs, theme.GetTheme("").TabInactiveStyle.Render(tabStr))
		}
	}
	return tabContainerStyle().Render(lipgloss.JoinHorizontal(lipgloss.Left, tabs...))
}

func tabContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("240"))
}
