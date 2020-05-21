package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	client "github.com/moov-io/identity/pkg/client"
)

func NewTestClient(handler http.Handler) *client.APIClient {
	mockHandler := MockClientHandler{
		handler: handler,
	}

	mockClient := &http.Client{
		Transport: &mockHandler,
	}

	config := client.NewConfiguration()
	config.HTTPClient = mockClient
	apiClient := client.NewAPIClient(config)

	return apiClient
}

type MockClientHandler struct {
	handler http.Handler
	ctx     *context.Context
}

func (h *MockClientHandler) RoundTrip(request *http.Request) (*http.Response, error) {
	writer := httptest.NewRecorder()

	h.handler.ServeHTTP(writer, request)
	return writer.Result(), nil
}
