package bitbucketclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/marcodellorto/bitbucket-runner-autoscaler/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ErrReader struct{}

func (e *ErrReader) Read([]byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestGetRunners(t *testing.T) {
	const (
		baseURL       = "https://baseurl.com"
		workspaceUUID = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		validResponse = `{
  "page": 1,
  "values": [
    {
      "uuid": "{b6d86128-0946-4fc8-90bc-6e501c0e869c}",
      "name": "test",
      "labels": [
        "self.hosted",
        "linux"
      ],
      "state": {
        "status": "UNREGISTERED",
        "updated_on": "2024-11-16T09:55:35.926932702Z",
        "cordoned": false
      },
      "created_on": "2024-11-16T09:55:35.926685218Z",
      "updated_on": "2024-11-16T09:55:35.926685218Z",
      "oauth_client": {
        "id": "randomid",
        "token_endpoint": "https://auth.atlassian.com/oauth/token",
        "audience": "api.atlassian.com"
      }
    }
  ],
  "size": 1,
  "pagelen": 1
}`
	)

	validGetRunnersResponse := &GetRunnersResponse{}
	_ = json.Unmarshal([]byte(validResponse), &validGetRunnersResponse)

	url := fmt.Sprintf("%s/internal/workspaces/%s/pipelines-config/runners?pagelen=%d", baseURL, workspaceUUID, Pagelen)

	tables := []struct {
		client           func() *mocks.HTTPClient
		expectedResponse *GetRunnersResponse
		expectedError    func() error
		name             string
	}{
		{
			name: "client returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{}, fmt.Errorf("something went wrong")).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				return fmt.Errorf("something went wrong")
			},
		},
		{
			name: "io.ReadAll returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{Body: io.NopCloser(&ErrReader{})}, nil).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				return fmt.Errorf("simulated read error")
			},
		},
		{
			name: "response status code is not 200",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader("{}"))}, nil).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				return fmt.Errorf("failed to fetch runners, status: 400, body: {}")
			},
		},
		{
			name: "json.Unmarshal returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{"))}, nil).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				mySyntaxError := &json.SyntaxError{}

				// Use reflection to set the unexported 'msg' field
				msgField := reflect.ValueOf(mySyntaxError).Elem().FieldByName("msg")
				if !msgField.IsValid() {
					panic("field 'msg' not found")
				}

				// Get an unsafe pointer to the 'msg' field
				msgPtr := unsafe.Pointer(msgField.UnsafeAddr())
				*(*string)(msgPtr) = "unexpected end of JSON input"

				// Set the exported 'Offset' field directly
				mySyntaxError.Offset = 1

				return mySyntaxError
			},
		},
		{
			name: "happy path",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(validResponse))}, nil).Once()

				return &m
			},
			expectedResponse: validGetRunnersResponse,
			expectedError: func() error {
				return nil
			},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			httpClient := table.client()

			c := New(httpClient, baseURL, workspaceUUID)

			resp, err := c.GetRunners()

			assert.Equal(t, table.expectedResponse, resp)
			assert.Equal(t, table.expectedError(), err)

			httpClient.AssertExpectations(t)
		})
	}
}

func TestGetRunner(t *testing.T) {
	const (
		baseURL              = "https://baseurl.com"
		workspaceUUID        = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		runnerUUID    string = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		validResponse        = `{
  "uuid": "{b6d86128-0946-4fc8-90bc-6e501c0e869c}",
  "name": "test",
  "labels": [
    "self.hosted",
    "linux"
  ],
  "state": {
    "status": "UNREGISTERED",
    "updated_on": "2024-11-16T09:55:35.926932702Z",
    "cordoned": false
  },
  "created_on": "2024-11-16T09:55:35.926685218Z",
  "updated_on": "2024-11-16T09:55:35.926685218Z",
  "oauth_client": {
    "id": "randomid",
    "token_endpoint": "https://auth.atlassian.com/oauth/token",
    "audience": "api.atlassian.com"
  }
}`
	)

	validGetRunnerResponse := &Runner{}
	_ = json.Unmarshal([]byte(validResponse), &validGetRunnerResponse)

	url := fmt.Sprintf("%s/internal/workspaces/%s/pipelines-config/runners/%s", baseURL, workspaceUUID, runnerUUID)

	tables := []struct {
		client           func() *mocks.HTTPClient
		expectedResponse *Runner
		expectedError    func() error
		name             string
	}{
		{
			name: "client returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{}, fmt.Errorf("something went wrong")).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				return fmt.Errorf("something went wrong")
			},
		},
		{
			name: "io.ReadAll returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{Body: io.NopCloser(&ErrReader{})}, nil).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				return fmt.Errorf("simulated read error")
			},
		},
		{
			name: "response status code is not 200",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader("{}"))}, nil).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				return fmt.Errorf("failed to fetch runner, status: 400, body: {}")
			},
		},
		{
			name: "json.Unmarshal returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{"))}, nil).Once()

				return &m
			},
			expectedResponse: nil,
			expectedError: func() error {
				mySyntaxError := &json.SyntaxError{}

				// Use reflection to set the unexported 'msg' field
				msgField := reflect.ValueOf(mySyntaxError).Elem().FieldByName("msg")
				if !msgField.IsValid() {
					panic("field 'msg' not found")
				}

				// Get an unsafe pointer to the 'msg' field
				msgPtr := unsafe.Pointer(msgField.UnsafeAddr())
				*(*string)(msgPtr) = "unexpected end of JSON input"

				// Set the exported 'Offset' field directly
				mySyntaxError.Offset = 1

				return mySyntaxError
			},
		},
		{
			name: "happy path",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Get", url).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(validResponse))}, nil).Once()

				return &m
			},
			expectedResponse: validGetRunnerResponse,
			expectedError: func() error {
				return nil
			},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			httpClient := table.client()

			c := New(httpClient, baseURL, workspaceUUID)

			resp, err := c.GetRunner(runnerUUID)

			assert.Equal(t, table.expectedResponse, resp)
			assert.Equal(t, table.expectedError(), err)

			httpClient.AssertExpectations(t)
		})
	}
}

func TestDeleteRunner(t *testing.T) {
	const (
		baseURL              = "https://baseurl.com"
		workspaceUUID        = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		runnerUUID    string = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
	)

	url := fmt.Sprintf("%s/internal/workspaces/%s/pipelines-config/runners/%s", baseURL, workspaceUUID, runnerUUID)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	tables := []struct {
		client        func() *mocks.HTTPClient
		expectedError func() error
		name          string
	}{
		{
			name: "client returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Do", req).Return(&http.Response{}, fmt.Errorf("something went wrong")).Once()

				return &m
			},
			expectedError: func() error {
				return fmt.Errorf("something went wrong")
			},
		},
		{
			name: "response status code is not 204",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Do", req).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader("{}"))}, nil).Once()

				return &m
			},
			expectedError: func() error {
				return fmt.Errorf("failed to delete runner, status: 400, body: {}")
			},
		},
		{
			name: "happy path",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Do", req).Return(&http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil).Once()

				return &m
			},
			expectedError: func() error {
				return nil
			},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			httpClient := table.client()

			c := New(httpClient, baseURL, workspaceUUID)

			err := c.DeleteRunner(runnerUUID)

			assert.Equal(t, table.expectedError(), err)

			httpClient.AssertExpectations(t)
		})
	}
}

func TestPostRunner(t *testing.T) {
	const (
		baseURL       string = "https://baseurl.com"
		workspaceUUID string = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		runnerUUID    string = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		validResponse        = `{
            "uuid": "{b6d86128-0946-4fc8-90bc-6e501c0e869c}",
            "name": "test",
            "labels": [
              "self.hosted",
              "linux"
            ],
            "state": {
              "status": "UNREGISTERED",
              "updated_on": "2024-11-16T09:55:35.926932702Z",
              "cordoned": false
            },
            "created_on": "2024-11-16T09:55:35.926685218Z",
            "updated_on": "2024-11-16T09:55:35.926685218Z",
            "oauth_client": {
              "id": "randomid",
              "token_endpoint": "https://auth.atlassian.com/oauth/token",
              "audience": "api.atlassian.com"
            }
          }`
	)

	validPostRunnerResponse := &Runner{}
	_ = json.Unmarshal([]byte(validResponse), &validPostRunnerResponse)

	url := fmt.Sprintf("%s/internal/workspaces/%s/pipelines-config/runners", baseURL, workspaceUUID)

	req := PostRunnerRequest{Name: "a name", Labels: []string{"a-label", "another-label"}}

	tables := []struct {
		client         func() *mocks.HTTPClient
		expectedResult *Runner
		expectedError  func() error
		name           string
	}{
		{
			name: "client returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Post", url, contentTypeApplicationJSON, mock.AnythingOfType("*bytes.Reader")).Return(&http.Response{}, fmt.Errorf("something went wrong")).Once()

				return &m
			},
			expectedResult: nil,
			expectedError: func() error {
				return fmt.Errorf("something went wrong")
			},
		},
		{
			name: "io.ReadAll returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Post", url, contentTypeApplicationJSON, mock.AnythingOfType("*bytes.Reader")).Return(&http.Response{Body: io.NopCloser(&ErrReader{})}, nil).Once()

				return &m
			},
			expectedResult: nil,
			expectedError: func() error {
				return fmt.Errorf("simulated read error")
			},
		},
		{
			name: "response status code is not 200 or 201",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Post", url, contentTypeApplicationJSON, mock.AnythingOfType("*bytes.Reader")).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader("{}"))}, nil).Once()

				return &m
			},
			expectedResult: nil,
			expectedError: func() error {
				return fmt.Errorf("failed to create runner, status: 400, body: {}")
			},
		},
		{
			name: "unable to unmarshal response body",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Post", url, contentTypeApplicationJSON, mock.AnythingOfType("*bytes.Reader")).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{"))}, nil).Once()

				return &m
			},
			expectedResult: nil,
			expectedError: func() error {
				return fmt.Errorf("error unmarshalling POST runner response: unexpected end of JSON input")
			},
		},
		{
			name: "happy path",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Post", url, contentTypeApplicationJSON, mock.AnythingOfType("*bytes.Reader")).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(validResponse))}, nil).Once()

				return &m
			},
			expectedResult: validPostRunnerResponse,
			expectedError: func() error {
				return nil
			},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			httpClient := table.client()

			c := New(httpClient, baseURL, workspaceUUID)

			res, err := c.PostRunner(req)

			assert.Equal(t, table.expectedError(), err)
			assert.Equal(t, table.expectedResult, res)

			httpClient.AssertExpectations(t)
		})
	}
}

func TestPutRunnerStatus(t *testing.T) {
	const (
		baseURL         string = "https://baseurl.com"
		workspaceUUID   string = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		runnerUUID      string = "e2f9c256-1843-4fd6-8456-2f8a1d94f8b5"
		runnerNewStatus string = "DISABLED"
	)

	tables := []struct {
		client         func() *mocks.HTTPClient
		expectedResult *Runner
		expectedError  func() error
		name           string
	}{
		{
			name: "client returns an error",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, fmt.Errorf("something went wrong")).Once()

				return &m
			},
			expectedError: func() error {
				return fmt.Errorf("failed to PUT runner status: something went wrong")
			},
		},
		{
			name: "response status code is not between 200 and 300",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader("{}"))}, nil).Once()

				return &m
			},
			expectedError: func() error {
				return fmt.Errorf("failed to update runner status, status: 400, body: {}")
			},
		},
		{
			name: "happy path",
			client: func() *mocks.HTTPClient {
				m := mocks.HTTPClient{}

				m.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil).Once()

				return &m
			},
			expectedError: func() error {
				return nil
			},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			httpClient := table.client()

			c := New(httpClient, baseURL, workspaceUUID)

			err := c.PutRunnerStatus(runnerUUID, runnerNewStatus)

			assert.Equal(t, table.expectedError(), err)

			httpClient.AssertExpectations(t)
		})
	}
}
