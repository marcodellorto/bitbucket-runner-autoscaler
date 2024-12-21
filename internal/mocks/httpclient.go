package mocks

import (
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HTTPClient struct {
	mock.Mock
}

func (m *HTTPClient) Get(url string) (resp *http.Response, err error) {
	args := m.Called(url)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *HTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, contentType, body)

	return args.Get(0).(*http.Response), args.Error(1)
}
func (m *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}
