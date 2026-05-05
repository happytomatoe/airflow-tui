package dag

import (
	"context"
	"strings"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/charmbracelet/bubbletea"
)

// dagsLoadedMsg is sent when DAGs are fetched (handled by main Model.Update)
type dagsLoadedMsg struct {
	dags []airflow.DAG
	err  error
}

// FetchDAGs returns a tea.Cmd that asynchronously loads DAGs from the client
func (m *DAGScreenModel) FetchDAGs() tea.Cmd {
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

// filterDags filters DAG list by case-insensitive ID match (from original model.go:886)
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
