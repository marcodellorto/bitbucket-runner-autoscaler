package bitbucketclient

import (
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
	}

	return
}
