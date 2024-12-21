package bitbucketclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/marcodellorto/bitbucket-runner-autoscaler/internal/domain/ports"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	GetAccessTokenContentTypeHeader string = "application/x-www-form-urlencoded"
	GetRunnersPath                  string = "/internal/workspaces/%s/pipelines-config/runners?pagelen=%d"
	GetRunnerPath                   string = "/internal/workspaces/%s/pipelines-config/runners/%s"
	DeleteRunnerPath                string = "/internal/workspaces/%s/pipelines-config/runners/%s"
	PostRunnerPath                  string = "/internal/workspaces/%s/pipelines-config/runners"
	PutRunnerStatusPath             string = "/internal/workspaces/%s/pipelines-config/runners/%s/state"
	contentTypeApplicationJSON      string = "application/json"
	Pagelen                         int    = 100
)

type BitbucketClient struct {
	client        ports.HTTPClient
	baseURL       string
	workspaceUUID string
	// TODO: add logger
}

func NewBitbucketClient(workspaceUUID, baseURL, accessTokenURL, clientID, clientSecret string) *BitbucketClient {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     accessTokenURL,
		Scopes:       []string{},
	}

	client := config.Client(context.Background())

	return New(client, baseURL, workspaceUUID)
}

func New(client ports.HTTPClient, baseURL, workspaceUUID string) *BitbucketClient {
	return &BitbucketClient{
		client:        client,
		baseURL:       baseURL,
		workspaceUUID: workspaceUUID,
	}
}

func (c *BitbucketClient) GetRunners() (response *GetRunnersResponse, err error) {
	resp, err := c.client.Get(c.baseURL + fmt.Sprintf(GetRunnersPath, c.workspaceUUID, Pagelen))
	if err != nil {
		// TODO: log error
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// TODO: Error reading response body: err
		return
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch runners, status: %d, body: %s", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		// TODO: log Error unmarshalling response: err
		return
	}

	return
}

func (c *BitbucketClient) GetRunner(runnerUUID string) (*Runner, error) {
	url := c.baseURL + fmt.Sprintf(GetRunnerPath, c.workspaceUUID, runnerUUID)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch runner, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var runner Runner
	if err := json.Unmarshal(body, &runner); err != nil {
		// TODO: log Error unmarshalling response: err
		return nil, err
	}

	return &runner, nil
}

func (c *BitbucketClient) DeleteRunner(runnerUUID string) (err error) {
	url := c.baseURL + fmt.Sprintf(DeleteRunnerPath, c.workspaceUUID, runnerUUID)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// No body expected on success; just check status code
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete runner, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *BitbucketClient) PostRunner(requestBody PostRunnerRequest) (*Runner, error) {
	url := c.baseURL + fmt.Sprintf(PostRunnerPath, c.workspaceUUID)

	bodyBytes, _ := json.Marshal(requestBody)

	resp, err := c.client.Post(url, contentTypeApplicationJSON, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create runner, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var runner Runner
	if err := json.Unmarshal(body, &runner); err != nil {
		return nil, fmt.Errorf("error unmarshalling POST runner response: %s", err.Error())
	}

	return &runner, nil
}

func (c *BitbucketClient) PutRunnerStatus(runnerUUID, newStatus string) error {
	url := c.baseURL + fmt.Sprintf(PutRunnerStatusPath, c.workspaceUUID, runnerUUID)

	requestBody := PutRunnerStatus{
		Status: newStatus,
	}

	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", contentTypeApplicationJSON)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to PUT runner status: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update runner status, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
