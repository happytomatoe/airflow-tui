package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type searchDialogModel struct {
	title string
	input textinput.Model
}

func renderSearchDialog(s *searchDialogModel, theme Theme) string {
	if s == nil {
		return ""
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.BorderColor).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)

	titleStyle := theme.TitleStyle.Align(lipgloss.Center).Padding(0, 0, 1, 0)

	content := titleStyle.Render(s.title) + "\n\n" +
		s.input.View() + "\n\n" +
		theme.MutedStyle.Align(lipgloss.Center).Render("[Enter] search  [Esc] cancel")

	return boxStyle.Render(content)
}
