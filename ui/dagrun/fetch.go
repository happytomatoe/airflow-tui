package dagrun

import (
	"context"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	tea "github.com/charmbracelet/bubbletea"
)

// FetchDAGRuns returns a tea.Cmd that fetches DAG runs for the given DAG ID
func (m *DAGRunScreenModel) FetchDAGRuns(dagID string) tea.Cmd {
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

// dagRunsLoadedMsg is the message returned after fetching DAG runs
type dagRunsLoadedMsg struct {
	dagID string
	runs  []airflow.DAGRun
	err   error
}
