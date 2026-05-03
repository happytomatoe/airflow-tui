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
	spinner    spinner.Model
	table      table.Model
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
			{Title: "DAG ID", Width: 48},
			{Title: "Owners", Width: 24},
			{Title: "Schedule", Width: 28},
			{Title: "Paused", Width: 8},
		}),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(true)
	styles.Selected = styles.Selected.Bold(true)
	t.SetStyles(styles)

	m := &Model{
		cfg:     cfg,
		spinner: s,
		table:   t,
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
		m.table.SetHeight(max(5, msg.Height-8))
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			if m.client == nil {
				return m, nil
			}
			m.loading = true
			m.err = nil
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
			m.table.SetRows(makeRows(msg.dags))
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
		body = m.table.View()
	}

	header := titleStyle.Render("DAG List")
	if m.activeName != "" {
		header += "\n" + mutedStyle.Render("Server: "+m.activeName)
	}

	footer := mutedStyle.Render("q quit  r refresh")
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
			strings.Join(derefStrings(dag.Owners), ", "),
			scheduleText(dag),
			boolText(dag.IsPaused),
		})
	}
	return rows
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

func derefStrings(v *[]string) []string {
	if v == nil {
		return nil
	}
	return *v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
