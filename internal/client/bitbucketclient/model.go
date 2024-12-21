package bitbucketclient

import "time"

type GetRunnersResponse struct {
	Values  []Runner `json:"values"`
	Page    int      `json:"page"`
	Size    int      `json:"size"`
	Pagelen int      `json:"pagelen"`
}

type State struct {
	UpdatedOn time.Time `json:"updated_on"`
	Status    string    `json:"status"`
	Cordoned  bool      `json:"cordoned"`
}

type OauthClient struct {
	ID            string `json:"id"`
	TokenEndpoint string `json:"token_endpoint"`
	Audience      string `json:"audience"`
}

type Runner struct {
	CreatedOn   time.Time   `json:"created_on"`
	UpdatedOn   time.Time   `json:"updated_on"`
	OauthClient OauthClient `json:"oauth_client"`
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	State       State       `json:"state"`
	Labels      []string    `json:"labels"`
}

type PostRunnerRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

type PutRunnerStatus struct {
	Status string `json:"status"`
}
