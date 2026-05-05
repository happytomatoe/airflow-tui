package task

import (
	"context"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	tea "github.com/charmbracelet/bubbletea"
)

type taskInstancesLoadedMsg struct {
	dagID    string
	dagRunID string
	tasks    []airflow.TaskInstance
	err      error
}

func FetchTaskInstances(client airflow.Client, dagID, dagRunID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tasks, err := client.GetTaskInstances(ctx, dagID, dagRunID)
		if err != nil {
			return taskInstancesLoadedMsg{dagID: dagID, dagRunID: dagRunID, err: err}
		}

		return taskInstancesLoadedMsg{dagID: dagID, dagRunID: dagRunID, tasks: tasks}
	}
}
