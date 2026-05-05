package config

import "github.com/charmbracelet/bubbles/table"

// Column width constants
const (
	emptyColWidth1 = 2
	emptyColWidth2 = 2
	statusColWidth = 6
	nameColWidth   = 20
	urlColWidth    = 60
)

// Columns defines the table columns for the Config screen
var Columns = []table.Column{
	{Title: "", Width: emptyColWidth1},
	{Title: "", Width: emptyColWidth2},
	{Title: "Status", Width: statusColWidth},
	{Title: "Name", Width: nameColWidth},
	{Title: "URL", Width: urlColWidth},
}
