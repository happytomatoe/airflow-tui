package log

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type logLoadedMsg struct {
	Content   string
	NextToken *string
	Err       error
}

func (m *LogScreenModel) FetchLogs(dagID, dagRunID, taskID string) tea.Cmd {
	return func() tea.Msg {
		if m.client == nil {
			return logLoadedMsg{Err: fmt.Errorf("no client available")}
		}

		timeoutChan := make(chan bool, 1)
		go func() {
			time.Sleep(8 * time.Second)
			timeoutChan <- true
		}()

		resultChan := make(chan logLoadedMsg, 1)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			content, nextToken, err := m.client.GetTaskLog(
				ctx,
				dagID,
				dagRunID,
				taskID,
				m.logTryNumber,
				true,
				m.logToken,
			)
			resultChan <- logLoadedMsg{
				Content:   content,
				NextToken: nextToken,
				Err:       err,
			}
		}()

		select {
		case result := <-resultChan:
			return result
		case <-timeoutChan:
			return logLoadedMsg{Err: fmt.Errorf("log loading timeout")}
		}
	}
}
