package config

import (
	"github.com/airflow-tui/airflow-tui/config"
	"github.com/charmbracelet/bubbles/table"
)

// NewConfigScreen creates a new ConfigScreenModel initialized with the given config
func NewConfigScreen(cfg config.Config) *ConfigScreenModel {
	t := table.New(
		table.WithColumns(Columns),
		table.WithFocused(true),
	)
	model := &ConfigScreenModel{
		configTable: t,
		cfg:         cfg,
	}
	model.refreshRows()
	return model
}

// SetConfig updates the configuration and refreshes the table rows
func (m *ConfigScreenModel) SetConfig(cfg config.Config) {
	m.cfg = cfg
	m.refreshRows()
}

// GetSelectedServer returns the currently selected server based on table cursor position
func (m *ConfigScreenModel) GetSelectedServer() *config.ServerConfig {
	cursor := m.configTable.Cursor()
	if cursor < 0 || cursor >= len(m.cfg.Servers) {
		return nil
	}
	return &m.cfg.Servers[cursor]
}

// refreshRows updates the table rows based on current config and state
func (m *ConfigScreenModel) refreshRows() {
	if len(m.cfg.Servers) == 0 {
		m.configTable.SetRows([]table.Row{})
		return
	}
	var rows []table.Row
	cursorPos := m.configTable.Cursor()
	for i, srv := range m.cfg.Servers {
		activeMark := " "
		cursorMark := " "
		status := "   "
		if srv.Name == m.activeName {
			activeMark = "*"
			if m.connected {
				status = "up"
			} else if m.err != nil {
				status = "down"
			} else {
				status = "..."
			}
		}
		if i == cursorPos {
			cursorMark = ">"
		}
		rows = append(rows, table.Row{activeMark, cursorMark, status, srv.Name, srv.URL})
	}
	m.configTable.SetRows(rows)
	// Restore cursor position (clamp to valid range)
	cursor := m.configTable.Cursor()
	if cursor < 0 || cursor >= len(rows) {
		cursor = 0
	}
	m.configTable.SetCursor(cursor)
}

// SetActiveServer updates the active server name and connection state
func (m *ConfigScreenModel) SetActiveServer(name string, connected bool, err error) {
	m.activeName = name
	m.connected = connected
	m.err = err
	m.refreshRows()
}

// SetSize updates the table dimensions based on terminal size
func (m *ConfigScreenModel) SetSize(width, height int) {
	// Config screen uses default table sizing
}
