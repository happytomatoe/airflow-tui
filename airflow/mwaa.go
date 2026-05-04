package airflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mwaa"
)

type MWAAClient struct {
	baseURL string
	client  *http.Client
	token   string
}

func NewMWAAClient(ctx context.Context, envName, profile, region string) (*MWAAClient, error) {
	token, baseURL, err := GetMWAACliToken(ctx, envName, profile, region)
	if err != nil {
		return nil, err
	}
	mwaa := &MWAAClient{
		baseURL: strings.TrimSuffix(baseURL, "/") + "/api/v2",
		client:  &http.Client{},
		token:   token,
	}
	return mwaa, nil
}

func GetMWAACliToken(ctx context.Context, envName, profile, region string) (string, string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return "", "", fmt.Errorf("failed to load AWS config: %w", err)
	}
	if region != "" {
		cfg.Region = region
	}
	client := mwaa.NewFromConfig(cfg)
	resp, err := client.CreateCliToken(ctx, &mwaa.CreateCliTokenInput{
		Name: awsString(envName),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get MWAA CLI token: %w", err)
	}
	hostname := awsStringValue(resp.WebServerHostname)
	return awsStringValue(resp.CliToken), "https://" + hostname, nil
}

func awsString(s string) *string {
	return &s
}

func awsStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (c *MWAAClient) GetDags(ctx context.Context) ([]DAG, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/dags", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var result DAGCollection
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Dags == nil {
		return nil, nil
	}
	return *result.Dags, nil
}

func (c *MWAAClient) GetDagStats(ctx context.Context, dagID string) ([]DagStatsCollectionItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *MWAAClient) GetDagRuns(ctx context.Context, dagID string) ([]DAGRun, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/dags/"+dagID+"/dagRuns", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var result DAGRunCollection
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.DagRuns == nil {
		return nil, nil
	}
	return *result.DagRuns, nil
}

func (c *MWAAClient) GetTaskInstances(ctx context.Context, dagID, dagRunID string) ([]TaskInstance, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/dags/"+dagID+"/dagRuns/"+dagRunID+"/taskInstances", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var result TaskInstanceCollection
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.TaskInstances == nil {
		return nil, nil
	}
	return *result.TaskInstances, nil
}

func (c *MWAAClient) GetTaskLog(ctx context.Context, dagID, dagRunID, taskID string, tryNumber int, fullContent bool, token *string) (string, *string, error) {
	return "", nil, fmt.Errorf("not implemented")
}

func (c *MWAAClient) ToggleDag(ctx context.Context, dagID string, paused bool) error {
	return fmt.Errorf("not implemented")
}

func (c *MWAAClient) TriggerDagRun(ctx context.Context, dagID string, conf string) error {
	return fmt.Errorf("not implemented")
}

func (c *MWAAClient) ClearDagRun(ctx context.Context, dagRunID, dagID string) error {
	return fmt.Errorf("not implemented")
}
