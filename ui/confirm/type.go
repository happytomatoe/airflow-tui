package confirm

import (
	"github.com/airflow-tui/airflow-tui/ui/theme"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

type ConfirmDialogModel struct {
	title    string
	message  string
	callback func() tea.Cmd
	visible   bool
}

func NewConfirmDialog() *ConfirmDialogModel {
	return &ConfirmDialogModel{}
}

func (m *ConfirmDialogModel) ShowConfirm(title, message string, callback func() tea.Cmd) {
	m.title = title
	m.message = message
	m.callback = callback
	m.visible = true
}

func (m *ConfirmDialogModel) Hide() {
	m.visible = false
}

func (m *ConfirmDialogModel) IsVisible() bool {
	return m.visible
}

func (m *ConfirmDialogModel) Confirm() tea.Cmd {
	m.visible = false
	if m.callback != nil {
		return m.callback()
	}
	return nil
}

// Update handles key messages for the confirm dialog
func (m *ConfirmDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "y":
			return m, m.Confirm()
		case "n":
			m.Hide()
			return m, nil
		}
	}
	return m, nil
}

// Init implements tea.Model
func (m *ConfirmDialogModel) Init() tea.Cmd {
	return nil
}

// View implements tea.Model interface
func (m *ConfirmDialogModel) View() string {
	if !m.visible {
		return ""
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.GetTheme("").BorderColor).
		Padding(1, 2).
		Width(60)

	titleStyle := theme.GetTheme("").TitleStyle.Padding(0, 0, 1, 0)
	messageStyle := theme.GetTheme("").MutedStyle.Padding(1, 0)
	buttonsStyle := lipgloss.NewStyle().
		Foreground(theme.GetTheme("").AccentColor)

	title := titleStyle.Render(m.title)
	message := messageStyle.Render(m.message)
	buttons := buttonsStyle.Render("[y] Yes  [n] No")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		message,
		"",
		buttons,
	)

	return boxStyle.Render(content)
}
