package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/airflow-tui/airflow-tui/config"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cfg        config.Config
	client     airflow.Client
	activeName string
	width      int
	height     int
	loading    bool
	err        error
	dags       []airflow.DAG
	filter     string
	searching  bool
	spinner    spinner.Model
	table      table.Model
	input      textinput.Model
}

type dagsLoadedMsg struct {
	dags []airflow.DAG
	err  error
}

func NewModel(cfg config.Config) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "DAG ID", Width: 64},
			{Title: "Schedule", Width: 32},
			{Title: "Paused", Width: 8},
		}),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(true)
	styles.Selected = styles.Selected.Bold(true)
	t.SetStyles(styles)

	input := textinput.New()
	input.Prompt = "/ "
	input.Placeholder = "find DAG by substring"
	input.CharLimit = 128

	m := &Model{
		cfg:     cfg,
		spinner: s,
		table:   t,
		input:   input,
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
		m.table.SetWidth(max(20, msg.Width-8))
		m.table.SetHeight(max(5, msg.Height-10))
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
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
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
		body = fmt.Sprintf("%s Loading DAGs from %s", m.spinner.View(), m.activeName)
	case m.err != nil:
		body = errorStyle.Render(fmt.Sprintf("Failed to load DAGs: %v", m.err)) + "\n" +
			mutedStyle.Render("Press r to retry.")
	default:
		body = m.searchView() + m.table.View()
	}

	header := titleStyle.Render("DAG List")
	if m.activeName != "" {
		header += "\n" + mutedStyle.Render("Server: "+m.activeName)
	}

	footer := mutedStyle.Render("q quit  r refresh  / search")
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
	m.table.SetRows(makeRows(filtered))
	if len(filtered) > 0 {
		m.table.SetCursor(0)
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
