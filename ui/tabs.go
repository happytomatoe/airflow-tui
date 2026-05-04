package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

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

func NewTabBar() TabBar {
	return TabBar{
		active: TabConfig,
		width:  0,
	}
}

func (t *TabBar) SetActive(tab Tab) {
	t.active = tab
}

func (t *TabBar) Active() Tab {
	return t.active
}

func (t *TabBar) SetWidth(width int) {
	t.width = width
}

func (t *TabBar) View() string {
	var tabs []string
	for i := range tabNames {
		tabStr := tabNames[i]
		if Tab(i) == t.active {
			tabs = append(tabs, activeTabStyle.Render(tabStr))
		} else {
			tabs = append(tabs, tabStyle.Render(tabStr))
		}
	}
	return tabContainerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, tabs...))
}

func (t *TabBar) HandleKey(msg string) bool {
	switch msg {
	case "1":
		t.active = TabConfig
		return true
	case "2":
		t.active = TabDags
		return true
	case "3":
		t.active = TabRuns
		return true
	case "4":
		t.active = TabTasks
		return true
	case "5":
		t.active = TabLogs
		return true
	case "left":
		if t.active > 0 {
			t.active--
		}
		return true
	case "right":
		if t.active < Tab(len(tabNames)-1) {
			t.active++
		}
		return true
	}
	return false
}

var (
	tabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("212")).
			Padding(0, 1).
			Bold(true)

	tabContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(lipgloss.Color("240"))
)

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
