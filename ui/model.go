package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/config"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type panel int

const (
	dagPanel panel = iota
	dagRunPanel
	taskPanel
	logPanel
)

type navigation struct {
	server  string
	dag     string
	dagRun  string
	task    string
	logType string
}

type Model struct {
	cfg           config.Config
	client        airflow.Client
	activeName    string
	width         int
	height        int
	panel         panel
	loading       bool
	err           error
	dags          []airflow.DAG
	dagRuns       []airflow.DAGRun
	filter        string
	searching     bool
	spinner       spinner.Model
	dagTable      table.Model
	runsTable     table.Model
	input         textinput.Model
	tabBar        TabBar
	showHelp      bool
	nav           navigation
	taskInstances []airflow.TaskInstance
	taskTable     table.Model
	logViewport   viewport.Model
	logContent    string
	logFollow     bool
	logTryNumber  int
	logToken      *string
	dagCode       string
	confirm       *confirmDialog
	visualMode    bool
	anchor        int
	selected      map[int]bool
	showGantt     bool
}

type dagsLoadedMsg struct {
	dags []airflow.DAG
	err  error
}

type dagRunsLoadedMsg struct {
	dagID string
	runs  []airflow.DAGRun
	err   error
}

type taskInstancesLoadedMsg struct {
	dagID    string
	dagRunID string
	tasks    []airflow.TaskInstance
	err      error
}

type logLoadedMsg struct {
	dagID     string
	dagRunID  string
	taskID    string
	content   string
	tryNumber int
	nextToken *string
	err       error
}

type tickLogMsg struct{}
type logTimeoutMsg struct{}

type dagCodeLoadedMsg struct {
	dagID string
	code  string
	err   error
}

type actionResultMsg struct {
	success bool
	message string
	err     error
}

func NewModel(cfg config.Config) *Model {
	// Initialize logger for this model
	logFile, err := os.OpenFile("ui_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create UI debug log file: %v\n", err)
	} else {
		log.SetOutput(logFile)
		log.SetLevel(log.DebugLevel)
		log.SetReportTimestamp(true)
		log.SetPrefix("UI-MODEL")
	}

	s := spinner.New()
	s.Spinner = spinner.Dot

	dagTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "DAG ID", Width: 64},
			{Title: "Schedule", Width: 32},
			{Title: "Paused", Width: 8},
		}),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	dagStyles := table.DefaultStyles()
	dagStyles.Header = dagStyles.Header.Bold(true)
	dagStyles.Selected = dagStyles.Selected.Bold(true)
	dagTable.SetStyles(dagStyles)

	runsTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "Run ID", Width: 40},
			{Title: "State", Width: 10},
			{Title: "Date", Width: 20},
			{Title: "Type", Width: 12},
		}),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	runsStyles := table.DefaultStyles()
	runsStyles.Header = runsStyles.Header.Bold(true)
	runsStyles.Selected = runsStyles.Selected.Bold(true)
	runsTable.SetStyles(runsStyles)

	input := textinput.New()
	input.Prompt = "/ "
	input.Placeholder = "find DAG by substring"
	input.CharLimit = 128

	taskTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "Task ID", Width: 40},
			{Title: "State", Width: 12},
			{Title: "Try #", Width: 6},
			{Title: "Duration", Width: 10},
		}),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	taskStyles := table.DefaultStyles()
	taskStyles.Header = taskStyles.Header.Bold(true)
	taskStyles.Selected = taskStyles.Selected.Bold(true)
	taskTable.SetStyles(taskStyles)

	m := &Model{
		cfg:          cfg,
		spinner:      s,
		panel:        dagPanel,
		dagTable:     dagTable,
		runsTable:    runsTable,
		input:        input,
		tabBar:       NewTabBar(),
		taskTable:    taskTable,
		logViewport:  viewport.New(80, 24),
		logTryNumber: 1,
		selected:     make(map[int]bool),
	}

	m.activeName, m.client = activeClient(cfg)
	if m.client != nil {
		m.loading = true
	}

	return m
}

func activeClient(cfg config.Config) (string, airflow.Client) {
	activeName := cfg.ActiveServer
	if activeName == "" && len(cfg.Servers) > 0 {
		activeName = cfg.Servers[0].Name
	}

	for _, server := range cfg.Servers {
		if server.Name != activeName {
			continue
		}

		var auth airflow.AuthProvider
		switch server.Auth.Type {
		case "basic":
			auth = &airflow.BasicAuth{Username: server.Auth.Username, Password: server.Auth.Password}
		case "token":
			auth = &airflow.StaticToken{Token: server.Auth.Token}
		case "mwaa":
			client, err := airflow.NewMWAASessionClient(context.Background(), server.URL, server.Auth.Profile, server.Auth.Region)
			if err != nil {
				return activeName, nil
			}
			return activeName, client
		}

		return activeName, airflow.NewAirflowApiClient(server.URL, auth)
	}

	return activeName, nil
}

func (m *Model) Init() tea.Cmd {
	if m.client == nil {
		return nil
	}
	return tea.Batch(m.spinner.Tick, m.loadDags())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dagTable.SetWidth(max(20, msg.Width-8))
		m.dagTable.SetHeight(max(5, msg.Height-10))
		m.runsTable.SetWidth(max(20, msg.Width-8))
		m.runsTable.SetHeight(max(5, msg.Height-10))
		m.logViewport.Width = msg.Width - 4
		m.logViewport.Height = msg.Height - 8
		return m, nil

	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
				m.input.Blur()
				m.input.SetValue("")
				m.filter = ""
				m.applyFilter()
				return m, nil
			case "enter":
				m.searching = false
				m.input.Blur()
				// Keep the current selection and filter
				return m, nil
			case "up", "down", "pgup", "pgdown", "home", "end":
				// Allow arrow keys and navigation keys to pass through to table navigation
				// Handle search input update first
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)

				// Apply filter in real-time only if it changed
				newFilter := strings.TrimSpace(m.input.Value())
				if newFilter != m.filter {
					m.filter = newFilter
					m.applyFilter()
				}

				// Then handle table navigation
				var tableCmd tea.Cmd
				switch m.panel {
				case dagPanel:
					m.dagTable, tableCmd = m.dagTable.Update(msg)
				case dagRunPanel:
					m.runsTable, tableCmd = m.runsTable.Update(msg)
				case taskPanel:
					m.taskTable, tableCmd = m.taskTable.Update(msg)
				case logPanel:
					m.logViewport, tableCmd = m.logViewport.Update(msg)
				}

				if tableCmd != nil {
					return m, tea.Batch(cmd, tableCmd)
				}
				return m, cmd
			}

			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)

			// Apply filter in real-time only if it changed
			newFilter := strings.TrimSpace(m.input.Value())
			if newFilter != m.filter {
				m.filter = newFilter
				m.applyFilter()
			}
			return m, cmd
		}

		// Handle escape key to clear loading state and stop spinner
		if msg.String() == "esc" && m.loading {
			m.loading = false
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.panel == dagPanel {
				row := m.dagTable.Cursor()
				if row >= 0 && row < len(m.dags) {
					dagID := derefString(m.dags[row].DagId)
					m.navigateToRuns(dagID)
					return m, m.loadDagRuns(dagID)
				}
			} else if m.panel == dagRunPanel {
				row := m.runsTable.Cursor()
				if row >= 0 && row < len(m.dagRuns) {
					dagRunID := derefString(m.dagRuns[row].DagRunId)
					dagID := m.nav.dag
					if dagID == "" {
						dagID = m.activeTabDagID()
					}
					m.navigateToTasks(dagID, dagRunID)
					return m, m.loadTaskInstances(dagID, dagRunID)
				}
			} else if m.panel == taskPanel {
				row := m.taskTable.Cursor()
				if row >= 0 && row < len(m.taskInstances) {
					taskID := derefString(m.taskInstances[row].TaskId)
					dagID := m.nav.dag
					dagRunID := m.nav.dagRun
					m.navigateToLogs(dagID, dagRunID, taskID)
					tryNum := 1
					if m.taskInstances[row].TryNumber != nil {
						tryNum = *m.taskInstances[row].TryNumber
					}
					return m, m.loadLog(dagID, dagRunID, taskID, tryNum, nil)
				}
			}
		case "esc":
			if m.showHelp {
				m.showHelp = false
			} else if m.panel == logPanel {
				m.panel = taskPanel
				m.tabBar.SetActive(TabTasks)
				m.logContent = ""
				m.logToken = nil
			} else if m.panel == taskPanel {
				m.panel = dagRunPanel
				m.tabBar.SetActive(TabRuns)
				m.taskInstances = nil
				m.nav.task = ""
			} else if m.panel == dagRunPanel {
				m.panel = dagPanel
				m.tabBar.SetActive(TabDags)
				m.dagRuns = nil
				m.nav.dag = ""
				m.nav.dagRun = ""
			}
			return m, nil
		case "/":
			m.searching = true
			m.input.SetValue(m.filter)
			m.input.CursorEnd()
			return m, m.input.Focus()
		case "r":
			if m.client == nil {
				return m, nil
			}
			m.loading = true
			m.err = nil
			m.filter = ""
			m.input.SetValue("")
			return m, m.loadDags()
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "p":
			return m.handleToggleDag()
		case "t":
			return m.handleTriggerDag()
		case "v":
			return m.handleViewDagCode()
		case "o":
			return m.handleOpenInBrowser()
		case "V":
			return m.handleVisualMode()
		case "c":
			return m.handleClearRuns()
		case "m":
			return m.handleMarkRuns()
		case "d":
			return m.handleShowGraph()
		case "F":
			if m.tabBar.Active() == TabLogs {
				m.logFollow = !m.logFollow
				if m.logFollow && m.logToken != nil {
					return m, m.loadLog(m.nav.dag, m.nav.dagRun, m.nav.task, m.logTryNumber, m.logToken)
				}
			}
			return m, nil
		case "g":
			if m.tabBar.Active() == TabTasks {
				m.showGantt = !m.showGantt
				return m, nil
			} else if m.tabBar.Active() == TabLogs {
				m.logViewport.GotoTop()
				return m, nil
			}
		case "G":
			if m.tabBar.Active() == TabLogs {
				m.logViewport.GotoBottom()
				return m, nil
			}
			return m, nil
		}

		if m.tabBar.Active() == TabLogs {
			switch msg.String() {
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				tryNum, _ := strconv.Atoi(msg.String())
				if tryNum > 0 && m.nav.task != "" {
					m.logTryNumber = tryNum
					m.logContent = ""
					m.logToken = nil
					m.loading = true
					return m, m.loadLog(m.nav.dag, m.nav.dagRun, m.nav.task, tryNum, nil)
				}
			}
		}

		if m.visualMode {
			switch msg.String() {
			case "j":
				return m.handleVisualMove(1)
			case "k":
				return m.handleVisualMove(-1)
			case " ":
				return m.handleVisualToggle()
			}
		}

		if m.confirm != nil && (msg.String() == "y" || msg.String() == "n" || msg.String() == "enter" || msg.String() == "esc") {
			return m.handleConfirm(msg.String())
		}

		if m.tabBar.HandleKey(msg.String()) {
			m.updatePanelFromTab()
			return m, nil
		}

	case spinner.TickMsg:
		// Only process spinner if actually loading and model is valid
		if !m.loading || m.client == nil {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case dagsLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.dags = msg.dags
			m.applyFilter()
		}
		return m, nil

	case dagRunsLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.dagRuns = sortDagRuns(msg.runs)
			m.updateRunsTable()
		}
		return m, nil

	case taskInstancesLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.taskInstances = msg.tasks
			m.updateTaskTable()
		}
		return m, nil

	case logLoadedMsg:
		// Always clear loading state, regardless of success or error
		m.loading = false

		// Handle the error
		if msg.err != nil {
			m.err = msg.err
			log.Debug("Log loading error", "error", msg.err)
			return m, nil
		}

		// Handle successful log loading
		m.err = nil
		if msg.content != "" {
			log.Debug("=== Log content received ===")
			log.Debug("Content", "content", msg.content, "length", len(msg.content))
			log.Debug("Metadata", "dagID", msg.dagID, "dagRunID", msg.dagRunID, "taskID", msg.taskID, "tryNumber", msg.tryNumber, "nextToken", msg.nextToken)
			log.Debug("=== End debug ===")

			m.logContent += msg.content
			m.logViewport.SetContent(m.logContent)
			if m.logFollow {
				m.logViewport.GotoBottom()
			}
		}
		m.logTryNumber = msg.tryNumber
		m.logToken = msg.nextToken

		// Set up next tick for follow mode
		var cmd tea.Cmd
		if m.logFollow && m.tabBar.Active() == TabLogs && m.logToken != nil {
			cmd = tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return tickLogMsg{}
			})
		}
		return m, cmd

	case tickLogMsg:
		if m.logFollow && m.tabBar.Active() == TabLogs && m.logToken != nil {
			return m, m.loadLog(m.nav.dag, m.nav.dagRun, m.nav.task, m.logTryNumber, m.logToken)
		}
		return m, nil

	case logTimeoutMsg:
		m.loading = false
		m.err = fmt.Errorf("Log loading timed out")
		return m, nil

	case dagCodeLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.dagCode = msg.code
		}
		return m, nil

	case actionResultMsg:
		m.loading = false
		if msg.success {
			m.err = nil
		} else {
			m.err = msg.err
		}
		return m, nil
	}

	var cmd tea.Cmd
	switch m.panel {
	case dagPanel:
		m.dagTable, cmd = m.dagTable.Update(msg)
	case dagRunPanel:
		m.runsTable, cmd = m.runsTable.Update(msg)
	case taskPanel:
		m.taskTable, cmd = m.taskTable.Update(msg)
	case logPanel:
		m.logViewport, cmd = m.logViewport.Update(msg)
	}
	return m, cmd
}

func (m *Model) View() string {
	var body string

	switch {
	case len(m.cfg.Servers) == 0:
		body = errorStyle.Render("No servers configured.") + "\n" +
			mutedStyle.Render("Use `airflow-tui config add` first.")
	case m.client == nil:
		body = errorStyle.Render("Active server could not be initialized.")
	case m.loading:
		body = fmt.Sprintf("%s %s", m.spinner.View(), m.loadingMessage())
	case m.err != nil:
		body = errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n" +
			mutedStyle.Render("Press r to retry.")
	default:
		body = m.bodyContent()
	}

	throbber := ""
	if m.loading {
		throbber = " " + m.spinner.View()
	}

	header := m.headerView(throbber)
	tabs := m.tabBar.View()
	footer := m.footerView()

	if m.showHelp {
		helpView := renderHelp(m.tabBar.Active())
		return appStyle.Render(header + "\n" + tabs + "\n\n" + body + "\n\n" + footer + "\n" + helpView)
	}
	return appStyle.Render(header + "\n" + tabs + "\n\n" + body + "\n\n" + footer)
}

func (m *Model) loadingMessage() string {
	switch m.tabBar.Active() {
	case TabConfig:
		return "Loading config"
	case TabDags:
		return fmt.Sprintf("Loading DAGs from %s", m.activeName)
	case TabRuns:
		return fmt.Sprintf("Loading DAG runs for %s", m.nav.dag)
	case TabTasks:
		return fmt.Sprintf("Loading tasks for %s", m.nav.dagRun)
	case TabLogs:
		return fmt.Sprintf("Loading logs for %s", m.nav.task)
	}
	return "Loading"
}

func (m *Model) headerView(throbber string) string {
	version := mutedStyle.Render("v0.1.0")
	headerRight := version + throbber

	breadcrumb := m.breadcrumbView()

	headerContent := lipgloss.JoinHorizontal(lipgloss.Right, headerRight)
	headerLine := headerStyle.Render(breadcrumb) + "\n" + headerContent

	return headerLine
}

func (m *Model) breadcrumbView() string {
	return m.breadcrumbViewWithSep(" > ")
}

func (m *Model) breadcrumbViewWithSep(sep string) string {
	var parts []string
	if m.activeName != "" {
		parts = append(parts, "config")
	}

	switch m.tabBar.Active() {
	case TabConfig:
	case TabDags:
		if m.nav.dag != "" {
			parts = append(parts, m.nav.dag)
		}
	case TabRuns:
		parts = append(parts, m.nav.dag)
	case TabTasks:
		parts = append(parts, m.nav.dag, m.nav.dagRun)
	case TabLogs:
		parts = append(parts, m.nav.dag, m.nav.dagRun, m.nav.task)
	}

	result := strings.Join(parts, sep)
	return lipgloss.JoinHorizontal(lipgloss.Left, result)
}

func (m *Model) bodyContent() string {
	switch m.tabBar.Active() {
	case TabConfig:
		return m.configView()
	case TabDags:
		return m.dagListView()
	case TabRuns:
		return m.dagRunsView()
	case TabTasks:
		return m.taskInstancesView()
	case TabLogs:
		return m.logsView()
	}
	return ""
}

func (m *Model) dagListView() string {
	return m.searchView() + m.dagTable.View()
}

func (m *Model) dagRunsView() string {
	if len(m.dagRuns) == 0 {
		return mutedStyle.Render("No DAG runs found")
	}
	return m.runsTable.View()
}

func (m *Model) taskInstancesView() string {
	if len(m.taskInstances) == 0 {
		return mutedStyle.Render("No task instances found")
	}
	if m.showGantt {
		gantt := NewGanttChart()
		gantt.SetTasks(m.taskInstances)
		gantt.SetWidth(m.width - 8)
		return gantt.View()
	}
	return m.taskTable.View()
}

func (m *Model) logsView() string {
	// SetContent is handled in logLoadedMsg to avoid resetting viewport scroll

	follow := ""
	if m.logFollow {
		follow = " [FOLLOW]"
	}
	tryNum := fmt.Sprintf("Try: %d", m.logTryNumber)
	header := mutedStyle.Render(tryNum+follow) + "\n"
	if m.logContent == "" {
		return header + mutedStyle.Render("No logs loaded")
	}
	scrollInfo := fmt.Sprintf(" %.0f%% ", m.logViewport.ScrollPercent()*100)
	return header + scrollInfo + "\n" + m.logViewport.View()
}

func (m *Model) configView() string {
	if len(m.cfg.Servers) == 0 {
		return mutedStyle.Render("No servers configured")
	}
	var rows []table.Row
	for _, srv := range m.cfg.Servers {
		activeMark := " "
		if srv.Name == m.activeName {
			activeMark = "*"
		}
		rows = append(rows, table.Row{activeMark, srv.Name, srv.URL})
	}
	configTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "", Width: 2},
			{Title: "Name", Width: 20},
			{Title: "URL", Width: 60},
		}),
		table.WithRows(rows),
		table.WithFocused(false),
	)
	return configTable.View()
}

func (m *Model) footerView() string {
	switch m.tabBar.Active() {
	case TabConfig:
		return mutedStyle.Render("q quit  1-5/←→ tabs  ? help")
	case TabDags:
		return mutedStyle.Render("q quit  r refresh  / search  1-5/←→ tabs  ? help  enter runs  p pause  t trigger")
	case TabRuns:
		return mutedStyle.Render("esc back  q quit  r refresh  1-5/←→ tabs  ? help  enter tasks  c clear  m mark")
	case TabTasks:
		return mutedStyle.Render("esc back  q quit  r refresh  1-5/←→ tabs  ? help  enter logs  c clear  m mark  o open")
	case TabLogs:
		return mutedStyle.Render("esc back  q quit  1-5/←→ tabs  ? help  F follow  gg top  G bottom  o open")
	}
	return ""
}

func (m *Model) loadDags() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		dags, err := m.client.GetDags(ctx)
		if err != nil {
			return dagsLoadedMsg{err: err}
		}

		return dagsLoadedMsg{dags: dags}
	}
}

func (m *Model) loadDagRuns(dagID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		runs, err := m.client.GetDagRuns(ctx, dagID)
		if err != nil {
			return dagRunsLoadedMsg{dagID: dagID, err: err}
		}

		return dagRunsLoadedMsg{dagID: dagID, runs: runs}
	}
}

func makeRows(dags []airflow.DAG) []table.Row {
	rows := make([]table.Row, 0, len(dags))
	for _, dag := range dags {
		rows = append(rows, table.Row{
			derefString(dag.DagId),
			scheduleText(dag),
			boolText(dag.IsPaused),
		})
	}
	return rows
}

func (m *Model) applyFilter() {
	filtered := filterDags(m.dags, m.filter)
	m.dagTable.SetRows(makeRows(filtered))
	if len(filtered) > 0 {
		m.dagTable.SetCursor(0)
	}
}

func (m *Model) searchView() string {
	if !m.searching && m.filter == "" {
		return ""
	}

	if m.searching {
		return m.input.View() + "\n\n"
	}

	return mutedStyle.Render("Filter: "+m.filter) + "\n\n"
}

func filterDags(dags []airflow.DAG, filter string) []airflow.DAG {
	if filter == "" {
		return dags
	}

	filter = strings.ToLower(filter)
	filtered := make([]airflow.DAG, 0, len(dags))
	for _, dag := range dags {
		if strings.Contains(strings.ToLower(derefString(dag.DagId)), filter) {
			filtered = append(filtered, dag)
		}
	}
	return filtered
}

func (m *Model) updateRunsTable() {
	rows := make([]table.Row, 0, len(m.dagRuns))
	for _, run := range m.dagRuns {
		state := "-"
		if run.State != nil {
			state = string(*run.State)
		}
		date := "-"
		if run.LogicalDate != nil {
			date = run.LogicalDate.Format("2006-01-02 15:04")
		}
		runType := "-"
		if run.RunType != nil {
			runType = string(*run.RunType)
		}
		rows = append(rows, table.Row{
			derefString(run.DagRunId),
			state,
			date,
			runType,
		})
	}
	m.runsTable.SetRows(rows)
	if len(rows) > 0 {
		m.runsTable.SetCursor(0)
	}
}

func sortDagRuns(runs []airflow.DAGRun) []airflow.DAGRun {
	sorted := append([]airflow.DAGRun(nil), runs...)
	sort.SliceStable(sorted, func(i, j int) bool {
		ti := time.Time{}
		if sorted[i].LogicalDate != nil {
			ti = *sorted[i].LogicalDate
		}
		tj := time.Time{}
		if sorted[j].LogicalDate != nil {
			tj = *sorted[j].LogicalDate
		}
		return ti.After(tj)
	})
	return sorted
}

func scheduleText(dag airflow.DAG) string {
	if dag.TimetableDescription != nil && *dag.TimetableDescription != "" {
		return *dag.TimetableDescription
	}
	return "-"
}

func boolText(v *bool) string {
	if v == nil {
		return "-"
	}
	if *v {
		return "yes"
	}
	return "no"
}

func derefString(v *string) string {
	if v == nil || *v == "" {
		return "-"
	}
	return *v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *Model) updatePanelFromTab() {
	switch m.tabBar.Active() {
	case TabConfig:
		m.panel = dagPanel
	case TabDags:
		m.panel = dagPanel
	case TabRuns:
		m.panel = dagRunPanel
	case TabTasks:
		m.panel = taskPanel
	case TabLogs:
		m.panel = logPanel
	}
}

func (m *Model) loadTaskInstances(dagID, dagRunID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tasks, err := m.client.GetTaskInstances(ctx, dagID, dagRunID)
		if err != nil {
			return taskInstancesLoadedMsg{dagID: dagID, dagRunID: dagRunID, err: err}
		}

		return taskInstancesLoadedMsg{dagID: dagID, dagRunID: dagRunID, tasks: tasks}
	}
}

func (m *Model) loadLog(dagID, dagRunID, taskID string, tryNumber int, token *string) tea.Cmd {
	return func() tea.Msg {
		// Don't proceed if not loading (prevents hanging)
		if !m.loading {
			return nil
		}

		// Set up a timeout to prevent hanging
		timeoutChan := make(chan bool, 1)
		go func() {
			time.Sleep(8 * time.Second)
			timeoutChan <- true
		}()

		resultChan := make(chan logLoadedMsg, 1)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			content, nextToken, err := m.client.GetTaskLog(ctx, dagID, dagRunID, taskID, tryNumber, true, token)
			resultChan <- logLoadedMsg{dagID: dagID, dagRunID: dagRunID, taskID: taskID, tryNumber: tryNumber, content: content, nextToken: nextToken, err: err}
		}()

		select {
		case result := <-resultChan:
			return result
		case <-timeoutChan:
			// Timeout occurred
			return logLoadedMsg{dagID: dagID, dagRunID: dagRunID, taskID: taskID, tryNumber: tryNumber, err: fmt.Errorf("Log loading timeout")}
		}
	}
}

func (m *Model) updateTaskTable() {
	rows := make([]table.Row, 0, len(m.taskInstances))
	for _, task := range m.taskInstances {
		state := "-"
		if task.State != nil {
			state = string(*task.State)
		}
		tryNum := "-"
		if task.TryNumber != nil {
			tryNum = fmt.Sprintf("%d", *task.TryNumber)
		}
		duration := "-"
		if task.Duration != nil {
			duration = fmt.Sprintf("%.2fs", *task.Duration)
		}
		rows = append(rows, table.Row{
			derefString(task.TaskId),
			state,
			tryNum,
			duration,
		})
	}
	m.taskTable.SetRows(rows)
	if len(rows) > 0 {
		m.taskTable.SetCursor(0)
	}
}

func (m *Model) navigateToRuns(dagID string) {
	m.nav.dag = dagID
	m.tabBar.SetActive(TabRuns)
	m.updatePanelFromTab()
	m.loading = true
	m.err = nil
	m.dagRuns = nil
}

func (m *Model) navigateToTasks(dagID, dagRunID string) {
	m.nav.dag = dagID
	m.nav.dagRun = dagRunID
	m.tabBar.SetActive(TabTasks)
	m.updatePanelFromTab()
	m.loading = true
	m.err = nil
	m.taskInstances = nil
}

func (m *Model) navigateToLogs(dagID, dagRunID, taskID string) {
	m.nav.dag = dagID
	m.nav.dagRun = dagRunID
	m.nav.task = taskID
	m.tabBar.SetActive(TabLogs)
	m.updatePanelFromTab()
	m.loading = true
	m.err = nil
	m.logContent = ""
	m.logTryNumber = 1
	m.logToken = nil
}

func (m *Model) activeTabDagID() string {
	switch m.tabBar.Active() {
	case TabDags:
		row := m.dagTable.Cursor()
		if row >= 0 && row < len(m.dags) {
			return derefString(m.dags[row].DagId)
		}
	case TabRuns:
		if m.nav.dag != "" {
			return m.nav.dag
		}
	}
	return ""
}

type confirmDialog struct {
	title   string
	message string
	action  func() tea.Cmd
}

func (m *Model) handleToggleDag() (tea.Model, tea.Cmd) {
	if m.tabBar.Active() != TabDags {
		return m, nil
	}
	row := m.dagTable.Cursor()
	if row < 0 || row >= len(m.dags) {
		return m, nil
	}
	dagID := derefString(m.dags[row].DagId)
	isPaused := m.dags[row].IsPaused != nil && *m.dags[row].IsPaused
	action := func() tea.Cmd {
		return func() tea.Msg {
			err := m.client.ToggleDag(context.Background(), dagID, !isPaused)
			if err != nil {
				return actionResultMsg{success: false, message: "Failed to toggle DAG", err: err}
			}
			action := "unpaused"
			if isPaused {
				action = "paused"
			}
			return actionResultMsg{success: true, message: fmt.Sprintf("DAG %s %s", dagID, action)}
		}
	}
	verb := "unpause"
	if isPaused {
		verb = "pause"
	}
	m.confirm = &confirmDialog{
		title:   "Toggle DAG",
		message: fmt.Sprintf("Are you sure you want to %s DAG %s?", verb, dagID),
		action:  action,
	}
	return m, nil
}

func (m *Model) handleTriggerDag() (tea.Model, tea.Cmd) {
	if m.tabBar.Active() != TabDags {
		return m, nil
	}
	row := m.dagTable.Cursor()
	if row < 0 || row >= len(m.dags) {
		return m, nil
	}
	dagID := derefString(m.dags[row].DagId)
	action := func() tea.Cmd {
		return func() tea.Msg {
			err := m.client.TriggerDagRun(context.Background(), dagID, "")
			if err != nil {
				return actionResultMsg{success: false, message: "Failed to trigger DAG", err: err}
			}
			return actionResultMsg{success: true, message: fmt.Sprintf("DAG %s triggered", dagID)}
		}
	}
	m.confirm = &confirmDialog{
		title:   "Trigger DAG",
		message: fmt.Sprintf("Are you sure you want to trigger DAG %s?", dagID),
		action:  action,
	}
	return m, nil
}

func (m *Model) handleViewDagCode() (tea.Model, tea.Cmd) {
	if m.tabBar.Active() != TabDags {
		return m, nil
	}
	row := m.dagTable.Cursor()
	if row < 0 || row >= len(m.dags) {
		return m, nil
	}
	dagID := derefString(m.dags[row].DagId)
	m.loading = true
	m.err = nil
	m.dagCode = ""
	return m, m.loadDagCode(dagID)
}

func (m *Model) handleOpenInBrowser() (tea.Model, tea.Cmd) {
	var dagID, dagRunID, taskID string

	switch m.tabBar.Active() {
	case TabDags:
		row := m.dagTable.Cursor()
		if row >= 0 && row < len(m.dags) {
			dagID = derefString(m.dags[row].DagId)
		}
	case TabRuns:
		dagID = m.nav.dag
		row := m.runsTable.Cursor()
		if row >= 0 && row < len(m.dagRuns) {
			dagRunID = derefString(m.dagRuns[row].DagRunId)
		}
	case TabTasks:
		dagID = m.nav.dag
		dagRunID = m.nav.dagRun
		row := m.taskTable.Cursor()
		if row >= 0 && row < len(m.taskInstances) {
			taskID = derefString(m.taskInstances[row].TaskId)
		}
	case TabLogs:
		dagID = m.nav.dag
		dagRunID = m.nav.dagRun
		taskID = m.nav.task
	}

	if dagID == "" {
		return m, nil
	}

	baseURL := ""
	for _, srv := range m.cfg.Servers {
		if srv.Name == m.activeName {
			baseURL = srv.URL
			break
		}
	}
	if baseURL == "" {
		return m, nil
	}

	url := baseURL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "dags/" + dagID
	if dagRunID != "" {
		url += "/graph?dag_run_id=" + dagRunID
		if taskID != "" {
			url += "&task_id=" + taskID
		}
	}

	m.err = nil
	go func() {
		cmd := exec.Command("xdg-open", url)
		cmd.Start()
	}()
	return m, nil
}

func runOpen(url string) error {
	cmd := exec.Command("xdg-open", url)
	return cmd.Start()
}

func (m *Model) handleConfirm(key string) (tea.Model, tea.Cmd) {
	if m.confirm == nil {
		return m, nil
	}

	if key == "y" || key == "enter" {
		m.confirm = nil
		m.loading = true
		return m, m.confirm.action()
	}

	m.confirm = nil
	return m, nil
}

func (m *Model) loadDagCode(dagID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		code, err := getDagCode(ctx, m.client, dagID)
		if err != nil {
			return dagCodeLoadedMsg{dagID: dagID, err: err}
		}

		return dagCodeLoadedMsg{dagID: dagID, code: code}
	}
}

func getDagCode(ctx context.Context, client airflow.Client, dagID string) (string, error) {
	type dagCodeResponse struct {
		DagCode string `json:"dag_code"`
	}
	return "", nil
}

func (m *Model) handleVisualMode() (tea.Model, tea.Cmd) {
	m.visualMode = !m.visualMode
	if m.visualMode {
		m.anchor = m.currentTable().Cursor()
		m.selected = make(map[int]bool)
		m.selected[m.anchor] = true
	} else {
		m.selected = make(map[int]bool)
		m.anchor = 0
	}
	return m, nil
}

func (m *Model) handleVisualMove(dir int) (tea.Model, tea.Cmd) {
	current := m.currentTable()
	newPos := current.Cursor() + dir
	maxPos := m.currentItemCount() - 1
	if newPos < 0 {
		newPos = 0
	}
	if newPos > maxPos {
		newPos = maxPos
	}
	current.SetCursor(newPos)

	start := m.anchor
	end := newPos
	if start > end {
		start, end = end, start
	}
	for i := start; i <= end; i++ {
		m.selected[i] = true
	}
	return m, nil
}

func (m *Model) handleVisualToggle() (tea.Model, tea.Cmd) {
	pos := m.currentTable().Cursor()
	if m.selected[pos] {
		delete(m.selected, pos)
	} else {
		m.selected[pos] = true
	}
	return m, nil
}

func (m *Model) currentTable() table.Model {
	switch m.panel {
	case dagPanel:
		return m.dagTable
	case dagRunPanel:
		return m.runsTable
	case taskPanel:
		return m.taskTable
	}
	return m.dagTable
}

func (m *Model) currentItemCount() int {
	switch m.panel {
	case dagPanel:
		return len(m.dags)
	case dagRunPanel:
		return len(m.dagRuns)
	case taskPanel:
		return len(m.taskInstances)
	}
	return 0
}

func (m *Model) selectionCount() int {
	return len(m.selected)
}

func (m *Model) visualIndicator() string {
	if !m.visualMode {
		return ""
	}
	count := m.selectionCount()
	if count == 0 {
		return mutedStyle.Render(" [VISUAL]")
	}
	return mutedStyle.Render(fmt.Sprintf(" [VISUAL %d selected]", count))
}

func (m *Model) handleClearRuns() (tea.Model, tea.Cmd) {
	if m.tabBar.Active() != TabRuns {
		return m, nil
	}

	selected := m.getSelectedRunIDs()
	if len(selected) == 0 {
		return m, nil
	}

	action := func() tea.Cmd {
		return func() tea.Msg {
			for _, run := range selected {
				err := m.client.ClearDagRun(context.Background(), run.runID, run.dagID)
				if err != nil {
					return actionResultMsg{success: false, message: fmt.Sprintf("Failed to clear %s", run.runID), err: err}
				}
			}
			return actionResultMsg{success: true, message: fmt.Sprintf("Cleared %d runs", len(selected))}
		}
	}

	m.confirm = &confirmDialog{
		title:   "Clear DAG Runs",
		message: fmt.Sprintf("Are you sure you want to clear %d selected DAG runs?", len(selected)),
		action:  action,
	}
	return m, nil
}

func (m *Model) handleMarkRuns() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) handleShowGraph() (tea.Model, tea.Cmd) {
	return m, nil
}

type selectedRun struct {
	dagID string
	runID string
}

func (m *Model) getSelectedRunIDs() []selectedRun {
	if len(m.selected) == 0 {
		row := m.runsTable.Cursor()
		if row >= 0 && row < len(m.dagRuns) {
			return []selectedRun{{dagID: m.nav.dag, runID: derefString(m.dagRuns[row].DagRunId)}}
		}
		return nil
	}

	var results []selectedRun
	for idx := range m.selected {
		if idx >= 0 && idx < len(m.dagRuns) {
			results = append(results, selectedRun{dagID: m.nav.dag, runID: derefString(m.dagRuns[idx].DagRunId)})
		}
	}
	return results
}
