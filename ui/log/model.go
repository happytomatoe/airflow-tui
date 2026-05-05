package log

import (
	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/charmbracelet/bubbles/viewport"
)

func NewLogScreen() *LogScreenModel {
	return &LogScreenModel{
		logViewport:  viewport.New(80, 24),
		logTryNumber: 1,
	}
}

func (m *LogScreenModel) SetLogs(content string, nextToken *string) {
	m.logContent = content
	m.logToken = nextToken
	m.logViewport.SetContent(content)
}

func (m *LogScreenModel) SetClient(client airflow.Client) {
	m.client = client
}

func (m *LogScreenModel) SetLogContent(content string) {
	m.logContent = content
	m.logViewport.SetContent(content)
}

func (m *LogScreenModel) ToggleFollow() {
	m.logFollow = !m.logFollow
}

func (m *LogScreenModel) NextTry() {
	m.logTryNumber++
}

func (m *LogScreenModel) PrevTry() {
	if m.logTryNumber > 1 {
		m.logTryNumber--
	}
}

// SetSize updates the viewport dimensions based on terminal size
func (m *LogScreenModel) SetSize(width, height int) {
	m.logViewport.Width = width - 8
	m.logViewport.Height = height - 10
}
