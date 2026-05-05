package tabs

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
