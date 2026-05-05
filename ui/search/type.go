package search

import (
	"github.com/airflow-tui/airflow-tui/ui/theme"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type SearchDialogModel struct {
	title string
	input textinput.Model
	visible bool
}

func NewSearchDialog() *SearchDialogModel {
	input := textinput.New()
	input.Prompt = "/ "
	input.Placeholder = "find DAG by substring"
	input.CharLimit = 128

	return &SearchDialogModel{
		input: input,
	}
}

func (m *SearchDialogModel) Show(title string) {
	m.title = title
	m.visible = true
	m.input.Focus()
}

func (m *SearchDialogModel) Hide() {
	m.visible = false
	m.input.Blur()
}

func (m *SearchDialogModel) IsVisible() bool {
	return m.visible
}

func (m *SearchDialogModel) GetFilter() string {
	return m.input.Value()
}

func (m *SearchDialogModel) SetFilter(filter string) {
	m.input.SetValue(filter)
}

func (m *SearchDialogModel) Update(msg interface{}) {
	if keyMsg, ok := msg.(interface{ String() string }); ok {
		m.input, _ = m.input.Update(keyMsg)
	}
}

// View implements tea.Model interface
func (m *SearchDialogModel) View() string {
	if !m.visible {
		return ""
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.GetTheme("").BorderColor).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)

	titleStyle := theme.GetTheme("").TitleStyle.
		Align(lipgloss.Center).
		Padding(0, 0, 1, 0)

	content := titleStyle.Render(m.title) + "\n\n" +
		m.input.View() + "\n\n" +
		theme.GetTheme("").MutedStyle.Align(lipgloss.Center).Render("[Enter] search  [Esc] cancel")

	return boxStyle.Render(content)
}
