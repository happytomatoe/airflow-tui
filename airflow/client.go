package airflow

import (
	"context"
)

type Client interface {
	GetDags(ctx context.Context) ([]DAG, error)
	GetDagStats(ctx context.Context, dagID string) ([]DagStatsCollectionItem, error)
	GetDagRuns(ctx context.Context, dagID string) ([]DAGRun, error)
	GetTaskInstances(ctx context.Context, dagID, dagRunID string) ([]TaskInstance, error)
	GetTaskLog(ctx context.Context, dagID, dagRunID, taskID string, tryNumber int, fullContent bool, token *string) (string, *string, error)
	ToggleDag(ctx context.Context, dagID string, paused bool) error
	TriggerDagRun(ctx context.Context, dagID string, conf string) error
	ClearDagRun(ctx context.Context, dagRunID, dagID string) error
}

func NewAirflowApiClient(baseURL string, auth AuthProvider) Client {
	client, err := newAirflowApiV1Client(baseURL, auth)
	if err != nil {
		return nil
	}
	return client
}
