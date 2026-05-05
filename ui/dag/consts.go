package dag

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Column width constants as specified
const (
	DagIDWidth    = 64
	ScheduleWidth = 32
	PausedWidth   = 8
)

// DagTableColumns defines the table structure for DAG listing
var DagTableColumns = []table.Column{
	{Title: lipgloss.NewStyle().Bold(true).Render("DAG ID"), Width: DagIDWidth},
	{Title: lipgloss.NewStyle().Bold(true).Render("Schedule"), Width: ScheduleWidth},
	{Title: lipgloss.NewStyle().Bold(true).Render("Paused"), Width: PausedWidth},
}
