package bitbucketclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/marcodellorto/bitbucket-runner-autoscaler/internal/domain/ports"
	"github.com/marcodellorto/bitbucket-runner-autoscaler/internal/infra/logger"
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
	logger        ports.Logger
}

func NewBitbucketClient(workspaceUUID, baseURL, accessTokenURL, clientID, clientSecret string) *BitbucketClient {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     accessTokenURL,
		Scopes:       []string{},
	}

	client := config.Client(context.Background())

	//debug level
	logger := logger.NewLogger(-4)

	return New(client, logger, baseURL, workspaceUUID)
}

func New(client ports.HTTPClient, logger ports.Logger, baseURL, workspaceUUID string) *BitbucketClient {
	return &BitbucketClient{
		client:        client,
		baseURL:       baseURL,
		workspaceUUID: workspaceUUID,
		logger:        logger,
	}
}

func (c *BitbucketClient) GetRunners() (response *GetRunnersResponse, err error) {
	resp, err := c.client.Get(c.baseURL + fmt.Sprintf(GetRunnersPath, c.workspaceUUID, Pagelen))
	if err != nil {
		c.logger.Error("failed to fetch runners",
			"error", err,
			"url", c.baseURL+fmt.Sprintf(GetRunnersPath, c.workspaceUUID, Pagelen),
		)
		return nil, err
	}

	c.logger.Debug("deferring the closure of the response", "debug", "defer resp.Body.Close()")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("failed to read the response body", "error", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch runners, status: %d, body: %s", resp.StatusCode, string(body))
	}

	c.logger.Info("Successfully got bitbucket runners", "info", resp.StatusCode)

	err = json.Unmarshal(body, &response)
	if err != nil {
		c.logger.Error("failed to unmarshal the body", "error", err)
		return
	}

	return
}

func (c *BitbucketClient) GetRunner(runnerUUID string) (*Runner, error) {
	url := c.baseURL + fmt.Sprintf(GetRunnerPath, c.workspaceUUID, runnerUUID)

	resp, err := c.client.Get(url)
	if err != nil {
		c.logger.Error("failed to fetch runners",
			"error", err,
			"url", c.baseURL+fmt.Sprintf(GetRunnerPath, c.workspaceUUID, runnerUUID),
		)
		return nil, err
	}

	c.logger.Debug("deferring the closure of the response", "debug", "defer resp.Body.Close()")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("failed to read the response body", "error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch runner, status: %d, body: %s", resp.StatusCode, string(body))
	}

	c.logger.Info("Successfully got bitbucket runner", "info", "runner-uuid="+runnerUUID)
	var runner Runner
	if err := json.Unmarshal(body, &runner); err != nil {
		c.logger.Error("failed to unmarshal the body", "error", err)
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
