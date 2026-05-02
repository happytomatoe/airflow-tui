package airflow

import (
	"context"
)

type Client interface {
	GetDags(ctx context.Context) ([]DAG, error)
	GetDagStats(ctx context.Context, dagID string) ([]DagStatsCollectionItem, error)
	GetDagRuns(ctx context.Context, dagID string) ([]DAGRun, error)
	GetTaskInstances(ctx context.Context, dagID, dagRunID string) ([]TaskInstance, error)
	GetTaskLog(ctx context.Context, dagID, dagRunID, taskID string, tryNumber int) (string, error)
	ToggleDag(ctx context.Context, dagID string, paused bool) error
	TriggerDagRun(ctx context.Context, dagID string, conf string) error
	ClearDagRun(ctx context.Context, dagRunID, dagID string) error
}

// NewAirflowApiV2Client creates a new v2 API client using the generated OpenAPI client
func NewAirflowApiV2Client(baseURL string, auth AuthProvider) Client {
	client, err := newAirflowApiV2Client(baseURL, auth)
	if err != nil {
		return nil
	}
	return client
}
