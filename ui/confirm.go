package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("212")).
			Padding(1, 2).
			Width(60)

	confirmTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")).
				Bold(true).
				Padding(0, 0, 1, 0)

	confirmMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Padding(1, 0)

	confirmButtonsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212"))
)

func renderConfirm(c *confirmDialog) string {
	if c == nil {
		return ""
	}

	title := confirmTitleStyle.Render(c.title)
	message := confirmMessageStyle.Render(c.message)
	buttons := confirmButtonsStyle.Render("[y] Yes  [n] No")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		message,
		"",
		buttons,
	)

	return confirmBoxStyle.Render(content)
}
