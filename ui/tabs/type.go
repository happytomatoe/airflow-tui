package tabs

import "github.com/charmbracelet/bubbles/key"

type Tab int

const (
	TabConfig Tab = iota
	TabDags
	TabRuns
	TabTasks
	TabLogs
)

var tabNames = [...]string{"Config", "DAGs", "Runs", "Tasks", "Logs"}

type TabBar struct {
	active Tab
	width  int
}

type TabKeyMap struct {
	SwitchTab key.Binding
	PrevTab   key.Binding
	NextTab   key.Binding
}

func (k TabKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SwitchTab, k.PrevTab, k.NextTab},
	}
}

func DefaultTabKeyMap() TabKeyMap {
	return TabKeyMap{
		SwitchTab: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5"),
			key.WithHelp("1-5", "switch tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "prev tab"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "next tab"),
		),
	}
}
