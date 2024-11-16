package bitbucketclient

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/marcodellorto/bitbucket-runner-autoscaler/internal/domain/ports"
)

const (
	GetRunnersPath string = "/%s/pipelines-config/runners?pagelen=%d"
	Pagelen        int    = 100
)

type BitbucketClient struct {
	client        ports.HTTPClient
	baseURL       string
	workspaceUUID string
	// TODO: add logger
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

	err = json.Unmarshal(body, &response)
	if err != nil {
		// TODO: log Error unmarshalling response: err
		return
	}

	return
}
