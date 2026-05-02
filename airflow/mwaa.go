package airflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type MWAAClient struct {
	baseURL string
	client  *http.Client
	token   string
}

func NewMWAAClient(ctx context.Context, baseURL, profile, region string) (*MWAAClient, error) {
	token, err := getMWAAToken(ctx, profile, region)
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

func getMWAAToken(ctx context.Context, profile, region string) (string, error) {
	return "", fmt.Errorf("MWAA auth requires AWS credentials - use aws-sdk-go-v2")
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
