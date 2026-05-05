package theme

import "github.com/charmbracelet/lipgloss"

// Base styles (used as defaults before theme is applied)
var (
	AppStyle = lipgloss.NewStyle().
		Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("212"))

	MutedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	LogViewStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Background(lipgloss.Color("236"))

	// Additional styles used by screens
	CellStyle = lipgloss.NewStyle()
	CursorStyle = lipgloss.NewStyle().
		Bold(true)
	PausedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))
	ActiveStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)
)

// Theme defines the styling for the entire TUI.
type Theme struct {
	AppStyle         lipgloss.Style
	TitleStyle       lipgloss.Style
	HeaderStyle      lipgloss.Style
	MutedStyle       lipgloss.Style
	ErrorStyle       lipgloss.Style
	LogViewStyle     lipgloss.Style
	SelectedStyle    lipgloss.Style
	TabActiveStyle   lipgloss.Style
	TabInactiveStyle lipgloss.Style
	BorderColor      lipgloss.Color
	AccentColor      lipgloss.Color
}

// GetTheme returns a theme by name.
func GetTheme(name string) Theme {
	switch name {
	case "dracula":
		return draculaTheme
	case "gruvbox":
		return gruvboxTheme
	case "solarized":
		return solarizedTheme
	case "nord":
		return nordTheme
	case "onedark":
		return oneDarkTheme
	default:
		return defaultTheme
	}
}
