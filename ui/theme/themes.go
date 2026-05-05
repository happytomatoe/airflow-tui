package theme

import "github.com/charmbracelet/lipgloss"

var defaultTheme = Theme{
	AppStyle: lipgloss.NewStyle().
		Padding(1, 2),
	TitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")),
	HeaderStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("212")),
	MutedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	ErrorStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true),
	LogViewStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Background(lipgloss.Color("236")),
	SelectedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")),
	TabActiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("212")),
	TabInactiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	BorderColor: lipgloss.Color("212"),
	AccentColor: lipgloss.Color("212"),
}

var draculaTheme = Theme{
	AppStyle: lipgloss.NewStyle().
		Padding(1, 2),
	TitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("189")),
	HeaderStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("189")),
	MutedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	ErrorStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("204")).
		Bold(true),
	LogViewStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("234")),
	SelectedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("189")),
	TabActiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("189")),
	TabInactiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	BorderColor: lipgloss.Color("189"),
	AccentColor: lipgloss.Color("189"),
}

var gruvboxTheme = Theme{
	AppStyle: lipgloss.NewStyle().
		Padding(1, 2),
	TitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("172")),
	HeaderStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("172")),
	MutedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	ErrorStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("167")).
		Bold(true),
	LogViewStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("235")),
	SelectedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("172")),
	TabActiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("172")),
	TabInactiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	BorderColor: lipgloss.Color("172"),
	AccentColor: lipgloss.Color("172"),
}

var solarizedTheme = Theme{
	AppStyle: lipgloss.NewStyle().
		Padding(1, 2),
	TitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("166")),
	HeaderStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("166")),
	MutedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")),
	ErrorStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("160")).
		Bold(true),
	LogViewStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("234")),
	SelectedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("166")),
	TabActiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("166")),
	TabInactiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")),
	BorderColor: lipgloss.Color("166"),
	AccentColor: lipgloss.Color("166"),
}

var nordTheme = Theme{
	AppStyle: lipgloss.NewStyle().
		Padding(1, 2),
	TitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("110")),
	HeaderStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("231")).
		Background(lipgloss.Color("110")),
	MutedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	ErrorStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")).
		Bold(true),
	LogViewStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("223")).
		Background(lipgloss.Color("235")),
	SelectedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("110")),
	TabActiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")).
		Background(lipgloss.Color("110")),
	TabInactiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	BorderColor: lipgloss.Color("110"),
	AccentColor: lipgloss.Color("110"),
}

var oneDarkTheme = Theme{
	AppStyle: lipgloss.NewStyle().
		Padding(1, 2),
	TitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("111")),
	HeaderStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("251")).
		Background(lipgloss.Color("111")),
	MutedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	ErrorStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("204")).
		Bold(true),
	LogViewStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("251")).
		Background(lipgloss.Color("234")),
	SelectedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("111")),
	TabActiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("251")).
		Background(lipgloss.Color("111")),
	TabInactiveStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
	BorderColor: lipgloss.Color("111"),
	AccentColor: lipgloss.Color("111"),
}
