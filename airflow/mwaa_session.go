package airflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

type MWAASessionClient struct {
	baseURL      string
	client       *http.Client
	cookies      []*http.Cookie
	profile      string
	region       string
	webserverURL string
}

func NewMWAASessionClient(ctx context.Context, webserverURL, profile, region string) (*MWAASessionClient, error) {
	if !strings.HasPrefix(webserverURL, "https://") {
		webserverURL = "https://" + webserverURL
	}
	webserverURL = strings.TrimSuffix(webserverURL, "/")

	c := &MWAASessionClient{
		baseURL:      webserverURL + "/api/v1",
		client:       &http.Client{Timeout: 30 * time.Second},
		profile:      profile,
		region:       region,
		webserverURL: webserverURL,
	}

	if err := c.login(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *MWAASessionClient) login(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(c.profile))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	if c.region != "" {
		cfg.Region = c.region
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve AWS credentials: %w", err)
	}

	loginURL := c.webserverURL + "/aws_mwaa/login"
	req, err := http.NewRequestWithContext(ctx, "GET", loginURL, nil)
	if err != nil {
		return err
	}

	signer := v4.NewSigner()
	if err := signer.SignHTTP(ctx, creds, req, bodyHash(""), c.region, "execute-api", time.Now()); err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("MWAA login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status during login: %d", resp.StatusCode)
	}

	c.cookies = resp.Cookies()
	if len(c.cookies) == 0 {
		return fmt.Errorf("no session cookie received from MWAA")
	}
	return nil
}

func bodyHash(body string) string {
	h := sha256.New()
	h.Write([]byte(body))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *MWAASessionClient) makeRequest(ctx context.Context, method, path string) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %w", err)
		}
		for _, cookie := range c.cookies {
			req.AddCookie(cookie)
		}
		return c.client.Do(req)
	}

	return resp, nil
}

func (c *MWAASessionClient) GetDags(ctx context.Context) ([]DAG, error) {
	resp, err := c.makeRequest(ctx, "GET", "/dags")
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

func (c *MWAASessionClient) GetDagStats(ctx context.Context, dagID string) ([]DagStatsCollectionItem, error) {
	resp, err := c.makeRequest(ctx, "GET", "/dagStats?dag_ids="+dagID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Dags *[]DagStatsCollectionItem `json:"dags,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Dags == nil {
		return nil, nil
	}
	return *result.Dags, nil
}

func (c *MWAASessionClient) GetDagRuns(ctx context.Context, dagID string) ([]DAGRun, error) {
	resp, err := c.makeRequest(ctx, "GET", "/dags/"+dagID+"/dagRuns")
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

func (c *MWAASessionClient) GetTaskInstances(ctx context.Context, dagID, dagRunID string) ([]TaskInstance, error) {
	resp, err := c.makeRequest(ctx, "GET", "/dags/"+dagID+"/dagRuns/"+dagRunID+"/taskInstances")
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

func (c *MWAASessionClient) GetTaskLog(ctx context.Context, dagID, dagRunID, taskID string, tryNumber int, fullContent bool, token *string) (string, *string, error) {
	url := "/dags/" + dagID + "/dagRuns/" + dagRunID + "/taskInstances/" + taskID + "/logs/" + fmt.Sprintf("%d", tryNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+url, nil)
	if err != nil {
		return "", nil, err
	}

	if fullContent {
		q := req.URL.Query()
		q.Set("full_content", "true")
		req.URL.RawQuery = q.Encode()
	}
	if token != nil {
		q := req.URL.Query()
		q.Set("token", *token)
		req.URL.RawQuery = q.Encode()
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var logResp struct {
		Content           string  `json:"content"`
		ContinuationToken *string `json:"continuation_token,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&logResp); err != nil {
		return "", nil, err
	}

	return logResp.Content, logResp.ContinuationToken, nil
}

func (c *MWAASessionClient) ToggleDag(ctx context.Context, dagID string, paused bool) error {
	body := map[string]bool{"is_paused": paused}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := c.baseURL + "/dags/" + dagID
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

func (c *MWAASessionClient) TriggerDagRun(ctx context.Context, dagID string, conf string) error {
	body := map[string]string{"conf": conf}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := c.baseURL + "/dags/" + dagID + "/dagRuns"
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

func (c *MWAASessionClient) ClearDagRun(ctx context.Context, dagRunID, dagID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+"/dags/"+dagID+"/dagRuns/"+dagRunID, nil)
	if err != nil {
		return err
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}
