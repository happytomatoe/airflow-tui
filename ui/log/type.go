package log

import (
	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/charmbracelet/bubbles/viewport"
)

type LogScreenModel struct {
	logViewport  viewport.Model
	logContent   string
	logFollow    bool
	logTryNumber int
	logToken     *string
	client       airflow.Client
}
