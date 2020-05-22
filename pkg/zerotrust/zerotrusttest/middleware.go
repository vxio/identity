package zerotrusttest

import (
	"context"
	"net/http"

	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/zerotrust"
)

// TestMiddleware - Handles injecting a session into a request for testing
type TestMiddleware struct {
	time    stime.TimeService
	session zerotrust.Session
}

// NewTestMiddleware - Generates a default Middleware that always injects the specified Session into the request
func NewTestMiddleware(time stime.TimeService, session zerotrust.Session) *TestMiddleware {
	return &TestMiddleware{
		time:    time,
		session: session,
	}
}

// Handler - Generates the handler you use to wrap the http routes
func (s *TestMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't really like using this map of any objects in the context for this, but it seems how its done.
		ctx := context.WithValue(r.Context(), zerotrust.SessionContextKey, &s.session)

		h.ServeHTTP(w, r.Clone(ctx))
	})
}