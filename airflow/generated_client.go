package airflow

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	generated "github.com/airflow-tui/airflow-tui/airflow/v2/generated"
)

type AirflowApiV2Client struct {
	client  *generated.ClientWithResponses
	baseURL string
	auth    AuthProvider
}

func newAirflowApiV2Client(baseURL string, auth AuthProvider) (*AirflowApiV2Client, error) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	serverURL := baseURL + "/api/v2/"

	c, err := generated.NewClientWithResponses(serverURL, generated.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		if auth != nil {
			return auth.ApplyAuth(req)
		}
		return nil
	}))
	if err != nil {
		return nil, err
	}

	return &AirflowApiV2Client{
		client:  c,
		baseURL: baseURL,
		auth:    auth,
	}, nil
}

func (c *AirflowApiV2Client) GetDags(ctx context.Context) ([]generated.DAG, error) {
	limit := 100
	params := generated.GetDagsParams{
		Limit: &limit,
	}
	resp, err := c.client.GetDagsWithResponse(ctx, &params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	if resp.JSON200.Dags == nil {
		return nil, nil
	}
	return *resp.JSON200.Dags, nil
}

func (c *AirflowApiV2Client) GetDagStats(ctx context.Context, dagID string) ([]generated.DagStatsCollectionItem, error) {
	params := &generated.GetDagStatsParams{
		DagIds: dagID,
	}
	resp, err := c.client.GetDagStatsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	if resp.JSON200.Dags == nil {
		return nil, nil
	}
	return *resp.JSON200.Dags, nil
}

func (c *AirflowApiV2Client) GetDagRuns(ctx context.Context, dagID string) ([]generated.DAGRun, error) {
	limit := 100
	params := generated.GetDagRunsParams{
		Limit: &limit,
	}
	resp, err := c.client.GetDagRunsWithResponse(ctx, generated.DAGID(dagID), &params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	if resp.JSON200.DagRuns == nil {
		return nil, nil
	}
	return *resp.JSON200.DagRuns, nil
}

func (c *AirflowApiV2Client) GetTaskInstances(ctx context.Context, dagID, dagRunID string) ([]generated.TaskInstance, error) {
	limit := 100
	params := generated.GetTaskInstancesParams{
		Limit: &limit,
	}
	resp, err := c.client.GetTaskInstancesWithResponse(ctx, generated.DAGID(dagID), generated.DAGRunID(dagRunID), &params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	if resp.JSON200.TaskInstances == nil {
		return nil, nil
	}
	return *resp.JSON200.TaskInstances, nil
}

func (c *AirflowApiV2Client) GetTaskLog(ctx context.Context, dagID, dagRunID, taskID string, tryNumber int) (string, error) {
	resp, err := c.client.GetLogWithResponse(ctx, generated.DAGID(dagID), generated.DAGRunID(dagRunID), generated.TaskID(taskID), generated.TaskTryNumber(tryNumber), nil)
	if err != nil {
		return "", err
	}
	if resp.JSON200 == nil {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	if resp.JSON200.Content == nil {
		return "", nil
	}
	return *resp.JSON200.Content, nil
}

func (c *AirflowApiV2Client) ToggleDag(ctx context.Context, dagID string, paused bool) error {
	params := generated.PatchDagParams{
		UpdateMask: &[]string{"is_paused"},
	}
	isPaused := paused
	body := generated.PatchDagJSONRequestBody{
		IsPaused: &isPaused,
	}
	resp, err := c.client.PatchDagWithResponse(ctx, generated.DAGID(dagID), &params, body)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	return nil
}

func (c *AirflowApiV2Client) TriggerDagRun(ctx context.Context, dagID string, conf string) error {
	dagRunID := "manual__" + dagID
	body := generated.PostDagRunJSONRequestBody{
		DagId:   &dagID,
		DagRunId: &dagRunID,
	}
	resp, err := c.client.PostDagRunWithResponse(ctx, generated.DAGID(dagID), body)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	return nil
}

func (c *AirflowApiV2Client) ClearDagRun(ctx context.Context, dagRunID, dagID string) error {
	resp, err := c.client.DeleteDagRunWithResponse(ctx, generated.DAGID(dagID), generated.DAGRunID(dagRunID))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	return nil
}
