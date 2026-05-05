package log

import (
	"fmt"
	"github.com/airflow-tui/airflow-tui/ui/theme"
)

func (m *LogScreenModel) Render() string {
	follow := ""
	if m.logFollow {
		follow = " [FOLLOW]"
	}
	tryNum := fmt.Sprintf("Try: %d", m.logTryNumber)
	header := theme.MutedStyle.Render(tryNum + follow) + "\n"

	if m.logContent == "" {
		return header + theme.MutedStyle.Render("No logs loaded")
	}

	scrollInfo := fmt.Sprintf(" %.0f%% ", m.logViewport.ScrollPercent()*100)
	return header + scrollInfo + "\n" + m.logViewport.View()
}
