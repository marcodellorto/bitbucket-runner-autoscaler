package mocks

import (
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
