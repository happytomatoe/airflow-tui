package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/config"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Panel int

const (
	PanelConfig Panel = iota
	PanelDAG
	PanelDAGRun
	PanelTaskInstance
	PanelLogs
)

type Model struct {
	panel  Panel
	width  int
	height int
	ticks  int

	cfg      config.Config
	selected int
	filter   string

	envData map[string]*EnvironmentData
	active  string

	error   string
	dagID   string
	runID   string
	taskID  string
	taskLog string

	tableDAG      table.Model
	tableDAGRun   table.Model
	tableTaskInst table.Model
}

type EnvironmentData struct {
	Dags     []airflow.DAG
	DagStats map[string][]airflow.DagStatsCollectionItem
	DagRuns  map[string][]airflow.DAGRun
	Tasks    map[string][]airflow.TaskInstance
	Client   airflow.Client
}

func NewModel(cfg config.Config) *Model {
	envData := make(map[string]*EnvironmentData)
	for _, server := range cfg.Servers {
		var auth airflow.AuthProvider
		switch server.Auth.Type {
		case "basic":
			auth = &airflow.BasicAuth{Username: server.Auth.Username, Password: server.Auth.Password}
		case "token":
			auth = &airflow.StaticToken{Token: server.Auth.Token}
		}

		var client airflow.Client
		client = airflow.NewAirflowApiV2Client(server.URL, auth)

		envData[server.Name] = &EnvironmentData{
			Client:   client,
			DagStats: make(map[string][]airflow.DagStatsCollectionItem),
			DagRuns:  make(map[string][]airflow.DAGRun),
			Tasks:    make(map[string][]airflow.TaskInstance),
		}
	}

	active := cfg.ActiveServer
	if active == "" && len(cfg.Servers) > 0 {
		active = cfg.Servers[0].Name
	}

	m := &Model{
		panel:   PanelDAG,
		envData: envData,
		active:  active,
		cfg:     cfg,
	}

	// Initialize DAG table with default columns
	dagColumns := []table.Column{
		{Title: "Active", Width: 10},
		{Title: "Name", Width: 30},
		{Title: "Owners", Width: 10},
		{Title: "Schedule", Width: 12},
		{Title: "Next Run", Width: 10},
		{Title: "Stats", Width: 15},
	}

	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true).Foreground(lipgloss.Color("86"))
	s.Selected = s.Selected.Bold(true).Foreground(lipgloss.Color("255")).Background(lipgloss.Color("56"))

	m.tableDAG = table.New(
		table.WithColumns(dagColumns),
		table.WithStyles(s),
		table.WithFocused(true),
	)

	// Initialize DAGRun table
	dagRunColumns := []table.Column{
		{Title: "Run ID", Width: 25},
		{Title: "Type", Width: 12},
		{Title: "State", Width: 12},
		{Title: "Start Date", Width: 20},
	}

	m.tableDAGRun = table.New(
		table.WithColumns(dagRunColumns),
		table.WithStyles(s),
		table.WithFocused(true),
	)

	// Initialize TaskInstance table
	taskColumns := []table.Column{
		{Title: "Task ID", Width: 30},
		{Title: "State", Width: 12},
		{Title: "Duration", Width: 15},
	}

	m.tableTaskInst = table.New(
		table.WithColumns(taskColumns),
		table.WithStyles(s),
		table.WithFocused(true),
	)

	if active != "" {
		m.loadDags()
	}

	return m
}

func (m *Model) loadDags() {
	if m.active == "" {
		return
	}
	env := m.envData[m.active]
	if env == nil || env.Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dags, err := env.Client.GetDags(ctx)
	if err != nil {
		m.error = err.Error()
		return
	}

	env.Dags = dags
	for _, dag := range dags {
		if dag.DagId == nil {
			continue
		}
		stats, err := env.Client.GetDagStats(ctx, *dag.DagId)
		if err != nil {
			continue
		}
		env.DagStats[*dag.DagId] = stats
	}

	m.updateTableRows()
}

func (m *Model) updateTableRows() {
	env := m.envData[m.active]
	if env == nil || len(env.Dags) == 0 {
		m.tableDAG.SetRows([]table.Row{})
		return
	}

	filtered := m.filterDags(env.Dags)
	pausedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))

	var rows []table.Row
	for _, dag := range filtered {
		// Active symbol
		activeSymbol := "◉"
		if dag.IsPaused != nil && *dag.IsPaused {
			activeSymbol = "◯"
		}

		// DAG ID
		dagID := ""
		if dag.DagId != nil {
			dagID = *dag.DagId
		}

		// Owners
		owners := ""
		if dag.Owners != nil && len(*dag.Owners) > 0 {
			owners = strings.Join(*dag.Owners, ",")
		}

		// Schedule
		schedule := "None"
		if dag.TimetableDescription != nil && *dag.TimetableDescription != "" {
			schedule = *dag.TimetableDescription
		}

		// Next Run
		nextRun := "None"
		if dag.NextDagrunCreateAfter != nil {
			nextRun = formatRemainingTime(*dag.NextDagrunCreateAfter)
		}

		// Stats
		stats := ""
		dagStats := env.DagStats
		if dagStats != nil {
			if statsItems, ok := dagStats[dagID]; ok {
				var statParts []string
				for _, statItem := range statsItems {
					if statItem.Stats != nil {
						for _, stat := range *statItem.Stats {
							if stat.State != nil && stat.Count != nil && *stat.Count > 0 {
								state := *stat.State
								count := *stat.Count
								var stateStyle lipgloss.Style
								switch state {
								case "running":
									stateStyle = greenStyle
								case "failed":
									stateStyle = redStyle
								case "success":
									stateStyle = greenStyle
								case "queued":
									stateStyle = yellowStyle
								default:
									stateStyle = pausedStyle
								}
								statParts = append(statParts, stateStyle.Render(fmt.Sprintf("%s:%d", state, count)))
							}
						}
					}
				}
				stats = strings.Join(statParts, " ")
			}
		}
		if stats == "" {
			stats = pausedStyle.Render("None")
		}

		rows = append(rows, table.Row{
			activeSymbol,
			dagID,
			owners,
			schedule,
			nextRun,
			stats,
		})
	}

	m.tableDAG.SetRows(rows)
	if len(rows) > 0 {
		if m.selected >= len(rows) {
			m.selected = len(rows) - 1
		}
		m.tableDAG.SetCursor(m.selected)
	}
}

func (m *Model) loadDagRuns(dagID string) {
	if m.active == "" || dagID == "" {
		return
	}
	env := m.envData[m.active]
	if env == nil || env.Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	runs, err := env.Client.GetDagRuns(ctx, dagID)
	if err != nil {
		m.error = err.Error()
		return
	}

	env.DagRuns[dagID] = runs
	m.updateDAGRunRows()
}

func (m *Model) updateDAGRunRows() {
	env := m.envData[m.active]
	if env == nil || m.dagID == "" {
		m.tableDAGRun.SetRows([]table.Row{})
		return
	}

	runs := env.DagRuns[m.dagID]
	if len(runs) == 0 {
		m.tableDAGRun.SetRows([]table.Row{})
		return
	}

	var rows []table.Row
	for _, run := range runs {
		runID := ""
		if run.DagRunId != nil {
			runID = *run.DagRunId
		}
		runType := ""
		if run.RunType != nil {
			runType = string(*run.RunType)
		}
		state := ""
		if run.State != nil {
			state = string(*run.State)
		}
		startDate := ""
		if run.StartDate != nil {
			startDate = run.StartDate.Format("2006-01-02 15:04:05")
		}

		rows = append(rows, table.Row{runID, runType, state, startDate})
	}

	m.tableDAGRun.SetRows(rows)
	if len(rows) > 0 {
		if m.selected >= len(rows) {
			m.selected = len(rows) - 1
		}
		m.tableDAGRun.SetCursor(m.selected)
	}
}

func (m *Model) loadTasks(dagID, runID string) {
	if m.active == "" || dagID == "" || runID == "" {
		return
	}
	env := m.envData[m.active]
	if env == nil || env.Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tasks, err := env.Client.GetTaskInstances(ctx, dagID, runID)
	if err != nil {
		m.error = err.Error()
		return
	}

	key := dagID + "/" + runID
	env.Tasks[key] = tasks
	m.updateTaskInstRows()
}

func (m *Model) updateTaskInstRows() {
	env := m.envData[m.active]
	if env == nil || m.dagID == "" || m.runID == "" {
		m.tableTaskInst.SetRows([]table.Row{})
		return
	}

	key := m.dagID + "/" + m.runID
	tasks := env.Tasks[key]
	if len(tasks) == 0 {
		m.tableTaskInst.SetRows([]table.Row{})
		return
	}

	var rows []table.Row
	for _, task := range tasks {
		taskID := ""
		if task.TaskId != nil {
			taskID = *task.TaskId
		}
		state := ""
		if task.State != nil {
			state = string(*task.State)
		}
		duration := ""
		if task.Duration != nil {
			duration = fmt.Sprintf("%s", *task.Duration)
		}

		rows = append(rows, table.Row{taskID, state, duration})
	}

	m.tableTaskInst.SetRows(rows)
	if len(rows) > 0 {
		if m.selected >= len(rows) {
			m.selected = len(rows) - 1
		}
		m.tableTaskInst.SetCursor(m.selected)
	}
}

func (m *Model) loadTaskLog(dagID, runID, taskID string) string {
	if m.active == "" || dagID == "" || runID == "" || taskID == "" {
		return ""
	}
	env := m.envData[m.active]
	if env == nil || env.Client == nil {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log, err := env.Client.GetTaskLog(ctx, dagID, runID, taskID, 1)
	if err != nil {
		m.error = err.Error()
		return ""
	}

	return log
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) nextPanel() {
	if m.panel < PanelLogs {
		m.panel++
	}
}

func (m *Model) previousPanel() {
	if m.panel > PanelConfig {
		m.panel--
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tableDAG.SetHeight(m.height - 8)
		m.tableDAGRun.SetHeight(m.height - 10)
		m.tableTaskInst.SetHeight(m.height - 10)
		// Update column widths on resize
		screenWidth := m.width - 2
		if screenWidth < 80 {
			screenWidth = 80
		}
		colWidths := calculateColWidths(screenWidth)
		columns := []table.Column{
			{Title: "Active", Width: colWidths[0]},
			{Title: "Name", Width: colWidths[1]},
			{Title: "Owners", Width: colWidths[2]},
			{Title: "Schedule", Width: colWidths[3]},
			{Title: "Next Run", Width: colWidths[4]},
			{Title: "Stats", Width: colWidths[5]},
		}
		m.tableDAG.SetColumns(columns)
		return m, nil
	case tickMsg:
		m.ticks++
		if m.panel == PanelDAG {
			m.loadDags()
		}
		return m, nil
	default:
		return m, nil
	}
}

type tickMsg time.Time

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		if m.panel == PanelDAG {
			return m, tea.Quit
		}
		m.previousPanel()
		return m, nil
	case "ctrl+c", "ctrl+d":
		return m, tea.Quit
	case "right", "l":
		m.nextPanel()
		if m.panel == PanelDAGRun && m.dagID != "" {
			m.loadDagRuns(m.dagID)
		}
		if m.panel == PanelTaskInstance && m.dagID != "" && m.runID != "" {
			m.loadTasks(m.dagID, m.runID)
		}
		return m, nil
	case "left", "h":
		m.previousPanel()
		return m, nil
	case "enter":
		if m.panel == PanelDAG && m.active != "" {
			env := m.envData[m.active]
			if env == nil {
				return m, nil
			}
			filtered := m.filterDags(env.Dags)
			sel := m.tableDAG.Cursor()
			if sel >= len(filtered) {
				return m, nil
			}
			dag := filtered[sel]
			if dag.DagId == nil {
				return m, nil
			}
			m.dagID = *dag.DagId
			m.panel = PanelDAGRun
			m.loadDagRuns(m.dagID)
		} else if m.panel == PanelDAGRun && m.active != "" && m.dagID != "" {
			env := m.envData[m.active]
			runs := env.DagRuns[m.dagID]
			sel := m.tableDAGRun.Cursor()
			if sel >= len(runs) {
				return m, nil
			}
			run := runs[sel]
			if run.DagRunId == nil {
				return m, nil
			}
			m.runID = *run.DagRunId
			m.panel = PanelTaskInstance
			m.loadTasks(m.dagID, m.runID)
		} else if m.panel == PanelTaskInstance && m.active != "" && m.dagID != "" && m.runID != "" {
			env := m.envData[m.active]
			key := m.dagID + "/" + m.runID
			tasks := env.Tasks[key]
			sel := m.tableTaskInst.Cursor()
			if sel >= len(tasks) {
				return m, nil
			}
			task := tasks[sel]
			if task.TaskId == nil {
				return m, nil
			}
			m.taskID = *task.TaskId
			m.taskLog = m.loadTaskLog(m.dagID, m.runID, m.taskID)
			m.panel = PanelLogs
		}
		return m, nil
	case "c":
		m.panel = PanelConfig
		return m, nil
	case "r":
		m.loadDags()
		return m, func() tea.Msg { return tickMsg(time.Now()) }
	case "up", "k":
		switch m.panel {
		case PanelDAG:
			m.tableDAG.MoveUp(1)
			m.selected = m.tableDAG.Cursor()
		case PanelDAGRun:
			m.tableDAGRun.MoveUp(1)
			m.selected = m.tableDAGRun.Cursor()
		case PanelTaskInstance:
			m.tableTaskInst.MoveUp(1)
			m.selected = m.tableTaskInst.Cursor()
		default:
			if m.selected > 0 {
				m.selected--
			}
		}
		return m, nil
	case "down", "j":
		switch m.panel {
		case PanelDAG:
			m.tableDAG.MoveDown(1)
			m.selected = m.tableDAG.Cursor()
		case PanelDAGRun:
			m.tableDAGRun.MoveDown(1)
			m.selected = m.tableDAGRun.Cursor()
		case PanelTaskInstance:
			m.tableTaskInst.MoveDown(1)
			m.selected = m.tableTaskInst.Cursor()
		default:
			m.selected++
		}
		return m, nil
	case "p":
		if m.panel == PanelDAG && m.active != "" {
			env := m.envData[m.active]
			if env == nil {
				return m, nil
			}
			filtered := m.filterDags(env.Dags)
			sel := m.tableDAG.Cursor()
			if sel >= len(filtered) {
				return m, nil
			}
			dag := filtered[sel]
			if dag.DagId == nil || dag.IsPaused == nil {
				return m, nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			env.Client.ToggleDag(ctx, *dag.DagId, !*dag.IsPaused)
			m.loadDags()
		}
		return m, nil
	case "t":
		if m.panel == PanelDAG && m.active != "" {
			env := m.envData[m.active]
			if env == nil {
				return m, nil
			}
			filtered := m.filterDags(env.Dags)
			sel := m.tableDAG.Cursor()
			if sel >= len(filtered) {
				return m, nil
			}
			dag := filtered[sel]
			if dag.DagId == nil {
				return m, nil
			}
			m.dagID = *dag.DagId
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			env.Client.TriggerDagRun(ctx, m.dagID, "{}")
			m.loadDagRuns(m.dagID)
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m *Model) View() string {
	header := m.renderHeader()
	main := m.renderPanel()
	footer := m.renderFooter()
	return lipgloss.JoinVertical(lipgloss.Left, header, main, footer)
}

func (m *Model) renderHeader() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220")).Render("Airflow TUI")
	active := lipgloss.NewStyle().Foreground(lipgloss.Color("green")).Render(" " + m.active)
	timeStr := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(time.Now().Format("15:04"))
	row := lipgloss.JoinHorizontal(lipgloss.Left, title, active, lipgloss.NewStyle().Width(m.width-30).Render(""), timeStr)
	return lipgloss.NewStyle().Height(1).Render(row)
}

func (m *Model) renderFooter() string {
	hints := "[q]quit [←→]nav [c]config [r]refresh [enter]view [p]pause [t]trigger"
	return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Height(1).Width(m.width).Render(hints)
}

func (m *Model) renderPanel() string {
	if m.error != "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("red")).Render(m.error)
	}

	switch m.panel {
	case PanelConfig:
		return m.renderConfig()
	case PanelDAG:
		return m.renderDAG()
	case PanelDAGRun:
		return m.renderDAGRun()
	case PanelTaskInstance:
		return m.renderTaskInstance()
	case PanelLogs:
		return m.renderLogs()
	default:
		return "Unknown panel"
	}
}

func (m *Model) renderConfig() string {
	title := lipgloss.NewStyle().Bold(true).Render("Server Config")
	var lines []string
	lines = append(lines, title)

	for _, server := range m.cfg.Servers {
		marker := "  "
		if server.Name == m.active {
			marker = "> "
		}
		lines = append(lines, marker+server.Name)
	}

	if len(m.cfg.Servers) == 0 {
		lines = append(lines, "", "No servers configured", "Use 'airflow-tui config add [name] [url]' to add one")
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(strings.Join(lines, "\n"))
}

func (m *Model) renderDAG() string {
	env := m.envData[m.active]
	if env == nil || len(env.Dags) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("No DAGs found")
	}

	// Update table dimensions
	screenWidth := m.width - 2
	if screenWidth < 80 {
		screenWidth = 80
	}

	colWidths := calculateColWidths(screenWidth)
	columns := []table.Column{
		{Title: "Active", Width: colWidths[0]},
		{Title: "Name", Width: colWidths[1]},
		{Title: "Owners", Width: colWidths[2]},
		{Title: "Schedule", Width: colWidths[3]},
		{Title: "Next Run", Width: colWidths[4]},
		{Title: "Stats", Width: colWidths[5]},
	}
	m.tableDAG.SetColumns(columns)
	m.tableDAG.SetHeight(m.height - 8)

	return m.tableDAG.View()
}

func formatRemainingTime(t time.Time) string {
	now := time.Now()
	duration := t.Sub(now)
	if duration < 0 {
		return "overdue"
	}
	seconds := int(duration.Seconds())
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%dh %dm", hours, minutes%60)
	}
	days := hours / 24
return fmt.Sprintf("%dd %dh", days, hours%24)
}

func calculateColWidths(screenWidth int) []int {
	minWidths := []int{6, 30, 10, 12, 10, 15}

	availableWidth := screenWidth - 3
	borderCount := 5
	availableWidth -= borderCount

	totalMin := 0
	for _, w := range minWidths {
		totalMin += w
	}

	if availableWidth > totalMin {
		minWidths[1] = availableWidth - (totalMin - 30)
	}

	return minWidths
}


func (m *Model) filterDags(dags []airflow.DAG) []airflow.DAG {
	if m.filter == "" {
		return dags
	}
	var filtered []airflow.DAG
	for _, dag := range dags {
		if dag.DagId != nil && strings.Contains(*dag.DagId, m.filter) {
			filtered = append(filtered, dag)
		}
	}
	return filtered
}

func (m *Model) renderDAGRun() string {
	title := lipgloss.NewStyle().Bold(true).Render("DAG Runs: " + m.dagID)
	env := m.envData[m.active]
	if env == nil {
		return lipgloss.JoinVertical(lipgloss.Left, title)
	}

	runs := env.DagRuns[m.dagID]
	if len(runs) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, title, lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("\nNo runs found"))
	}

	// Update table dimensions
	m.tableDAGRun.SetHeight(m.height - 10)

	return lipgloss.JoinVertical(lipgloss.Left, title, "", m.tableDAGRun.View())
}

func (m *Model) renderTaskInstance() string {
	title := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Tasks: %s / %s", m.dagID, m.runID))
	env := m.envData[m.active]
	if env == nil {
		return lipgloss.JoinVertical(lipgloss.Left, title)
	}

	key := m.dagID + "/" + m.runID
	tasks := env.Tasks[key]
	if len(tasks) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, title, lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("\nNo tasks found"))
	}

	// Update table dimensions
	m.tableTaskInst.SetHeight(m.height - 10)

	return lipgloss.JoinVertical(lipgloss.Left, title, "", m.tableTaskInst.View())
}

func (m *Model) renderLogs() string {
	title := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Logs: %s / %s / %s", m.dagID, m.runID, m.taskID))
	if m.taskLog == "" {
		return lipgloss.JoinVertical(lipgloss.Left, title, lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("\nNo logs available"))
	}

	logContent := m.taskLog
	if len(logContent) > m.height-5 {
		logContent = logContent[:m.height-5]
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(logContent))
}
