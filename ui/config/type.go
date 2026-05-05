package config

import (
	"github.com/airflow-tui/airflow-tui/config"
	"github.com/charmbracelet/bubbles/table"
)

// ConfigScreenModel represents the state of the Config screen
type ConfigScreenModel struct {
	configTable table.Model
	cfg         config.Config
	activeName  string
	connected   bool
	err         error
}
