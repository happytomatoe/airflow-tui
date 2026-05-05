package tabs

func NewTabBar() *TabBar {
	return &TabBar{
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
