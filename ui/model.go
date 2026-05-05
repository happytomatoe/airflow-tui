package ui

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/config"
	"github.com/airflow-tui/airflow-tui/ui/dag"
	"github.com/airflow-tui/airflow-tui/ui/dagrun"
	"github.com/airflow-tui/airflow-tui/ui/task"
	"github.com/airflow-tui/airflow-tui/ui/log"
	uiconfig "github.com/airflow-tui/airflow-tui/ui/config"
	"github.com/airflow-tui/airflow-tui/ui/confirm"
	"github.com/airflow-tui/airflow-tui/ui/help"
	"github.com/airflow-tui/airflow-tui/ui/search"
	"github.com/airflow-tui/airflow-tui/ui/gantt"
	"github.com/airflow-tui/airflow-tui/ui/tabs"
	"github.com/airflow-tui/airflow-tui/ui/nav"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	charmlog "github.com/charmbracelet/log"
)

// Async message types (kept in main model)
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

// Model is the main application model that orchestrates screens and dialogs
type Model struct {
	cfg           config.Config
	client        airflow.Client
	activeName    string
	width         int
	height        int
	loading       bool
	err           error
	connected     bool
	spinner       spinner.Model

	// Screen models (pointers to refactored screen packages)
	dagScreen     *dag.DAGScreenModel
	dagRunScreen  *dagrun.DAGRunScreenModel
	taskScreen    *task.TaskScreenModel
	logScreen     *log.LogScreenModel
	configScreen  *uiconfig.ConfigScreenModel

	// Dialog models (pointers to refactored dialog packages)
	confirmDialog *confirm.ConfirmDialogModel
	helpDialog    *help.HelpDialogModel
	searchDialog  *search.SearchDialogModel
	ganttDialog   *gantt.GanttDialogModel

	// Tab bar for navigation
	tabBar        *tabs.TabBar

	// Navigation state
	nav           nav.Navigation
}

// Helper functions (kept in main model)
func derefString(v *string) string {
	if v == nil || *v == "" {
		return "-"
	}
	return *v
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

func scheduleText(dag airflow.DAG) string {
	if dag.TimetableDescription != nil && *dag.TimetableDescription != "" {
		return *dag.TimetableDescription
	}
	return "-"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// activeClient initializes the Airflow client for the active server
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
			client, err := airflow.NewMWAAClient(context.Background(), server.URL, server.Auth.Profile, server.Auth.Region)
			if err != nil {
				return activeName, nil
			}
			return activeName, client
		}

		return activeName, airflow.NewAirflowApiClient(server.URL, auth)
	}

	return activeName, nil
}

// NewModel initializes the main model with all screen and dialog models
func NewModel(cfg config.Config) *Model {
	// Initialize debug logger
	logFile, err := os.OpenFile("ui_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create UI debug log file: %v\n", err)
	} else {
		charmlog.SetOutput(logFile)
		charmlog.SetLevel(charmlog.DebugLevel)
		charmlog.SetReportTimestamp(true)
		charmlog.SetPrefix("UI-MODEL")
	}

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Initialize screen models via their constructors
	dagScreen := dag.NewDAGScreen()
	dagRunScreen := dagrun.NewDAGRunScreen()
	taskScreen := task.NewTaskScreen()
	logScreen := log.NewLogScreen()
	configScreen := uiconfig.NewConfigScreen(cfg)

	// Initialize dialog models via their constructors
	confirmDialog := confirm.NewConfirmDialog()
	helpDialog := help.NewHelpDialog()
	searchDialog := search.NewSearchDialog()
	ganttDialog := gantt.NewGanttDialog()

	// Initialize tab bar
	tabBar := tabs.NewTabBar()

	m := &Model{
		cfg:           cfg,
		spinner:       s,
		dagScreen:     dagScreen,
		dagRunScreen:  dagRunScreen,
		taskScreen:    taskScreen,
		logScreen:     logScreen,
		configScreen:  configScreen,
		confirmDialog: confirmDialog,
		helpDialog:    helpDialog,
		searchDialog:  searchDialog,
		ganttDialog:   ganttDialog,
		tabBar:        tabBar,
	}

	m.activeName, m.client = activeClient(cfg)
	if m.client != nil {
		m.loading = true
	}

	return m
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	if m.client == nil {
		return nil
	}
	return tea.Batch(m.spinner.Tick, m.loadDags())
}

// Update implements tea.Model - handles ALL keys and async messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to all screens
		m.dagScreen.SetSize(msg.Width, msg.Height)
		m.dagRunScreen.SetSize(msg.Width, msg.Height)
		m.taskScreen.SetSize(msg.Width, msg.Height)
		m.logScreen.SetSize(msg.Width, msg.Height)
		m.configScreen.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Handle dialog keys first if any dialog is visible
		if m.confirmDialog.IsVisible() {
			switch msg.String() {
			case "y":
				return m, m.confirmDialog.Confirm()
			case "n":
				m.confirmDialog.Hide()
				return m, nil
			}
			return m, nil
		}
		if m.helpDialog.IsVisible() {
			if msg.String() == "?" {
				m.helpDialog.Hide()
			}
			return m, nil
		}
		if m.searchDialog.IsVisible() {
			// Handle search dialog keys
			if msg.String() == "enter" {
				m.searchDialog.Hide()
				return m, nil
			}
			if msg.String() == "esc" {
				m.searchDialog.Hide()
				return m, nil
			}
			// Pass key to search input
			m.searchDialog.Update(msg)
			return m, nil
		}
		if m.ganttDialog.IsVisible() {
			if msg.String() == "esc" || msg.String() == "g" {
				m.ganttDialog.Hide()
			}
			return m, nil
		}

		// Handle keys directly based on active tab
		switch m.tabBar.Active() {
		case tabs.TabDags:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "r":
				if m.client != nil {
					m.loading = true
					return m, m.loadDags()
				}
			case "/":
				m.searchDialog.Show("Filter DAGs")
			}
		case tabs.TabRuns:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "esc":
				m.tabBar.SetActive(tabs.TabDags)
			case "r":
				if m.client != nil {
					m.loading = true
					return m, m.loadDagRuns(m.nav.Dag)
				}
			}
		case tabs.TabTasks:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "esc":
				m.tabBar.SetActive(tabs.TabRuns)
			case "r":
				if m.client != nil {
					m.loading = true
					return m, m.loadTaskInstances(m.nav.Dag, m.nav.DagRun)
				}
			}
		case tabs.TabLogs:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "esc":
				m.tabBar.SetActive(tabs.TabTasks)
			}
		case tabs.TabConfig:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			}
		}

	case dagsLoadedMsg:
		m.loading = false
		m.err = msg.err
		m.connected = msg.err == nil
		if msg.err == nil {
			m.dagScreen.SetDAGs(msg.dags)
		}
		return m, nil

	case dagRunsLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.dagRunScreen.SetDAGRuns(msg.runs)
		}
		return m, nil

	case taskInstancesLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.taskScreen.SetTasks(msg.tasks)
		}
		return m, nil

	case logLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.logScreen.SetLogs(msg.content, msg.nextToken)
		}
		return m, nil

	case tickLogMsg:
		// Log following feature - simplified for now
		return m, nil

	case spinner.TickMsg:
		if !m.loading || m.client == nil {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model - checks dialog visibility first, then renders active screen
func (m *Model) View() string {
	// Render visible dialog instead of screen content
	if m.confirmDialog.IsVisible() {
		return m.confirmDialog.View()
	}
	if m.helpDialog.IsVisible() {
		return m.helpDialog.Render(m.tabBar.Active())
	}
	if m.searchDialog.IsVisible() {
		return m.searchDialog.View()
	}
	if m.ganttDialog.IsVisible() {
		return m.ganttDialog.View()
	}

	// Render active screen
		switch m.tabBar.Active() {
		case tabs.TabDags:
			return m.headerView() + "\n" + m.tabBar.View() + "\n\n" + m.dagScreen.Render() + "\n\n" + m.footerView()
		case tabs.TabRuns:
			return m.headerView() + "\n" + m.tabBar.View() + "\n\n" + m.dagRunScreen.Render() + "\n\n" + m.footerView()
		case tabs.TabTasks:
			return m.headerView() + "\n" + m.tabBar.View() + "\n\n" + m.taskScreen.Render() + "\n\n" + m.footerView()
		case tabs.TabLogs:
			return m.headerView() + "\n" + m.tabBar.View() + "\n\n" + m.logScreen.Render() + "\n\n" + m.footerView()
		case tabs.TabConfig:
			return m.headerView() + "\n" + m.tabBar.View() + "\n\n" + m.configScreen.Render() + "\n\n" + m.footerView()
		}

	return ""
}

// loadDags returns a command that loads DAGs from the API
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

// loadDagRuns returns a command that loads DAG runs for a specific DAG
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

// loadTaskInstances returns a command that loads task instances for a DAG run
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

// loadLog returns a command that loads log content for a task
func (m *Model) loadLog(dagID, dagRunID, taskID string, tryNumber int, token *string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		content, nextToken, err := m.client.GetTaskLog(ctx, dagID, dagRunID, taskID, tryNumber, true, token)
		return logLoadedMsg{
			dagID:     dagID,
			dagRunID:  dagRunID,
			taskID:    taskID,
			content:   content,
			tryNumber: tryNumber,
			nextToken: nextToken,
			err:       err,
		}
	}
}

// headerView renders the application header
func (m *Model) headerView() string {
	// Simplified header - just show connection status
	return m.connectionStatus()
}

// footerView renders the application footer
func (m *Model) footerView() string {
	// Simplified footer
	return "Press ? for help"
}

// connectionStatus renders the connection status indicator
func (m *Model) connectionStatus() string {
	status := "disconnected"
	if m.connected {
		status = "connected"
	}
	if m.loading {
		status = "connecting..."
	}
	return "● " + status
}
