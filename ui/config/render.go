package config

import (
	"github.com/airflow-tui/airflow-tui/ui/theme"
)

// Render generates the string representation of the Config screen
func (m *ConfigScreenModel) Render() string {
	if len(m.cfg.Servers) == 0 {
		return theme.GetTheme("").MutedStyle.Render("No servers configured")
	}
	m.refreshRows()
	return m.configTable.View()
}
