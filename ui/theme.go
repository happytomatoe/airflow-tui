package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	StyleNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("236"))

	StyleSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Background(lipgloss.Color("239"))

	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("220"))

	StyleSubtle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	StyleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("red"))

	StyleSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("green"))

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("220"))

	StyleFooter = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))
)
