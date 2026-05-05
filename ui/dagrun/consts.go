package dagrun

import "github.com/charmbracelet/bubbles/table"

// Columns defines the DAG run table columns with specified widths
var Columns = []table.Column{
	{Title: "Run ID", Width: 40},
	{Title: "State", Width: 10},
	{Title: "Date", Width: 20},
	{Title: "Type", Width: 12},
}
