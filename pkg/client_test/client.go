package client_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	client "github.com/moov-io/identity/pkg/client"
)

func NewTestClient(handler http.Handler) *client.APIClient {
	mockHandler := MockClientHandler{
		handler: handler,
	}

	mockClient := &http.Client{

		// Mock handler that sends the request to the handler passed in and returns the response without a server
		// middleman.
		Transport: &mockHandler,

		// Disables following redirects for testing.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},

		Timeout: 15 * time.Minute,
	}

	config := client.NewConfiguration()
	config.HTTPClient = mockClient
	apiClient := client.NewAPIClient(config)

	return apiClient
}

type MockClientHandler struct {
	handler http.Handler
}

func (h *MockClientHandler) RoundTrip(request *http.Request) (*http.Response, error) {
	writer := httptest.NewRecorder()

	h.handler.ServeHTTP(writer, request)
	return writer.Result(), nil
}
