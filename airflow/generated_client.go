package airflow

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	generated "github.com/airflow-tui/airflow-tui/airflow/generated"
	"github.com/charmbracelet/log"
)

// V1LogRegex is a compiled regex for parsing V1 log content (Python tuple format)
// This matches tuples like ('hostname', 'log content') or ('hostname', "log content")
// where the second element can use either single or double quotes.
var V1LogRegex = regexp.MustCompile(`\(\s*'((?:\\.|[^'])*)'\s*,\s*(?:"((?:\\.|[^"])*)"|'((?:\\.|[^'])*)')\s*\)`)

// parseV1LogContent parses V1 log content from Python tuple format to plain text.
// V1 Airflow logs come as serialized Python tuples: [('host', 'log line\nmore')]
// This extracts the log text, joins multiple tuples, and expands escaped newlines.
func parseV1LogContent(content string) string {
	log.Debug("parseV1LogContent input", "content", content)
	
	fragments := make([]string, 0)
	
	// Find all matches and extract the log content (second element of each tuple)
	matches := V1LogRegex.FindAllStringSubmatch(content, -1)
	log.Debug("Found regex matches", "count", len(matches))
	
	for i, match := range matches {
		log.Debug("Match", "index", i, "match", match)
		// Second element can be in group 2 (double quotes) or group 3 (single quotes)
		var fragment string
		if len(match) >= 3 {
			if match[2] != "" {
				fragment = match[2] // Double quoted
			} else if match[3] != "" {
				fragment = match[3] // Single quoted
			}
		}
		
		if fragment != "" {
			// Replace escaped newlines with actual newlines
			fragment = strings.ReplaceAll(fragment, "\\n", "\n")
			fragments = append(fragments, fragment)
		}
	}
	
	if len(fragments) == 0 {
		log.Debug("No V1 format matches found, returning original content")
		return content
	}
	
	// Join all fragments with newlines
	result := strings.Join(fragments, "\n")
	log.Debug("parseV1LogContent result", "result", result)
	return result
}

type AirflowApiClient struct {
	client  *generated.ClientWithResponses
	baseURL string
	auth    AuthProvider
}

func newAirflowApiV1Client(baseURL string, auth AuthProvider) (*AirflowApiClient, error) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	serverURL := baseURL + "/api/v1/"

	c, err := generated.NewClientWithResponses(serverURL, generated.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		if auth != nil {
			return auth.ApplyAuth(req)
		}
		return nil
	}))
	if err != nil {
		return nil, err
	}

	return &AirflowApiClient{
		client:  c,
		baseURL: baseURL,
		auth:    auth,
	}, nil
}

func (c *AirflowApiClient) GetDags(ctx context.Context) ([]generated.DAG, error) {
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

func (c *AirflowApiClient) GetDagStats(ctx context.Context, dagID string) ([]generated.DagStatsCollectionItem, error) {
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

func (c *AirflowApiClient) GetDagRuns(ctx context.Context, dagID string) ([]generated.DAGRun, error) {
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

func (c *AirflowApiClient) GetTaskInstances(ctx context.Context, dagID, dagRunID string) ([]generated.TaskInstance, error) {
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

func (c *AirflowApiClient) GetTaskLog(ctx context.Context, dagID, dagRunID, taskID string, tryNumber int, fullContent bool, token *string) (string, *string, error) {
	params := &generated.GetLogParams{}
	if fullContent {
		fc := generated.FullContent(true)
		params.FullContent = &fc
	}
	if token != nil {
		ct := generated.ContinuationToken(*token)
		params.Token = &ct
	}
	resp, err := c.client.GetLogWithResponse(ctx, generated.DAGID(dagID), generated.DAGRunID(dagRunID), generated.TaskID(taskID), generated.TaskTryNumber(tryNumber), params)
	if err != nil {
		return "", nil, err
	}
	if resp.JSON200 != nil && resp.JSON200.Content != nil {
		content := *resp.JSON200.Content
		log.Debug("Raw log content from API", "content", content)
		
		// Parse V1 log format (Python tuple)
		content = parseV1LogContent(content)
		log.Debug("After parseV1LogContent", "content", content)
		
		lines := strings.Split(content, "\n")
		if len(lines) > 1 && !strings.Contains(lines[0], " ") {
			log.Debug("Skipping hostname line", "hostname", lines[0])
			content = strings.Join(lines[1:], "\n")
			log.Debug("After hostname skip", "content", content)
		}
		
		var nextToken *string
		if resp.JSON200.ContinuationToken != nil {
			nt := *resp.JSON200.ContinuationToken
			nextToken = &nt
		}
		return content, nextToken, nil
	}
	if len(resp.Body) > 0 {
		return string(resp.Body), nil, nil
	}
	return "", nil, fmt.Errorf("unexpected status: %d", resp.StatusCode())
}

func (c *AirflowApiClient) ToggleDag(ctx context.Context, dagID string, paused bool) error {
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

func (c *AirflowApiClient) TriggerDagRun(ctx context.Context, dagID string, conf string) error {
	dagRunID := "manual__" + dagID
	body := generated.PostDagRunJSONRequestBody{
		DagId:    &dagID,
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

func (c *AirflowApiClient) ClearDagRun(ctx context.Context, dagRunID, dagID string) error {
	resp, err := c.client.DeleteDagRunWithResponse(ctx, generated.DAGID(dagID), generated.DAGRunID(dagRunID))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
	return nil
}
