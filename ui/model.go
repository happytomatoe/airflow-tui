package ui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/config"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type panel int

const (
	dagPanel panel = iota
	dagRunPanel
)

type Model struct {
	cfg        config.Config
	client     airflow.Client
	activeName string
	width      int
	height     int
	panel      panel
	loading    bool
	err        error
	dags       []airflow.DAG
	dagRuns    []airflow.DAGRun
	filter     string
	searching  bool
	spinner    spinner.Model
	dagTable   table.Model
	runsTable  table.Model
	input      textinput.Model
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

func NewModel(cfg config.Config) *Model {
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

	m := &Model{
		cfg:        cfg,
		spinner:    s,
		panel:     dagPanel,
		dagTable:  dagTable,
		runsTable: runsTable,
		input:     input,
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
		return m, nil

	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			case "enter":
				m.filter = strings.TrimSpace(m.input.Value())
				m.searching = false
				m.input.Blur()
				m.applyFilter()
				return m, nil
			}

			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.panel == dagPanel {
				row := m.dagTable.Cursor()
				if row >= 0 && row < len(m.dags) {
					dagID := derefString(m.dags[row].DagId)
					m.panel = dagRunPanel
					m.loading = true
					m.err = nil
					m.dagRuns = nil
					return m, m.loadDagRuns(dagID)
				}
			}
		case "esc":
			if m.panel == dagRunPanel {
				m.panel = dagPanel
				m.dagRuns = nil
				return m, nil
			}
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
		}

	case spinner.TickMsg:
		if !m.loading {
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
	}

	var cmd tea.Cmd
	if m.panel == dagPanel {
		m.dagTable, cmd = m.dagTable.Update(msg)
	} else if m.panel == dagRunPanel {
		m.runsTable, cmd = m.runsTable.Update(msg)
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
		if m.panel == dagRunPanel {
			body = fmt.Sprintf("%s Loading DAG runs", m.spinner.View())
		} else {
			body = fmt.Sprintf("%s Loading DAGs from %s", m.spinner.View(), m.activeName)
		}
	case m.err != nil:
		if m.panel == dagRunPanel {
			body = errorStyle.Render(fmt.Sprintf("Failed to load DAG runs: %v", m.err)) + "\n" +
				mutedStyle.Render("Press r to retry.")
		} else {
			body = errorStyle.Render(fmt.Sprintf("Failed to load DAGs: %v", m.err)) + "\n" +
				mutedStyle.Render("Press r to retry.")
		}
	default:
		if m.panel == dagRunPanel {
			body = m.runsTable.View()
		} else {
			body = m.searchView() + m.dagTable.View()
		}
	}

	header := titleStyle.Render("DAG List")
	if m.panel == dagRunPanel {
		header = titleStyle.Render("DAG Runs")
	}
	if m.activeName != "" {
		header += "\n" + mutedStyle.Render("Server: "+m.activeName)
	}

	footer := mutedStyle.Render("q quit  r refresh  / search")
	if m.panel == dagRunPanel {
		footer = mutedStyle.Render("esc back  q quit  r refresh")
	}
	return appStyle.Render(header + "\n\n" + body + "\n\n" + footer)
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
