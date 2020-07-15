package authn

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/tumbler/pkg/jwe"
)

// Middleware - Handles authenticating a request came from the authn services
type Middleware struct {
	log        logging.Logger
	time       stime.TimeService
	jweService jwe.JWEService
}

// NewMiddleware - Generates a default AuthnMiddleware for use with authenticating a request came from the authn services
func NewMiddleware(log logging.Logger, time stime.TimeService, jweService jwe.JWEService) (*Middleware, error) {
	return &Middleware{
		log:        log,
		time:       time,
		jweService: jweService,
	}, nil
}

// Handler - Generates the handler you use to wrap the http routes
func (s *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.FromRequest(r)
		if err != nil {
			s.log.Error().LogError("Session error", err)
			w.WriteHeader(404)
			return
		}

		// Don't really like using this map of any objects in the context for this, but it seems how its done.
		ctx := context.WithValue(r.Context(), LoginSessionContextKey, session)

		h.ServeHTTP(w, r.Clone(ctx))
	})
}

// FromRequest - Pulls out authenticationd details from the Request and calls Parse.
func (s *Middleware) FromRequest(r *http.Request) (*LoginSession, error) {

	// This is kind of a hack but we don't want the origin checking on a callback. So we're just going to set the correct value.
	target, err := jwe.GetTarget(r)
	r.Header.Set("Origin", target)
	if err != nil {
		return nil, err
	}

	cookie, err := GetAuthnCookie(r)
	if err != nil {
		return nil, s.log.Error().LogError("No session cookie found", err)
	}

	session := LoginSession{}
	_, err = s.jweService.ParseEncrypted(r, cookie.Value, &session)
	if err != nil {
		return nil, s.log.Error().LogErrorF("Session token parse failure - %w", err)
	}

	_, err = uuid.Parse(session.CredentialID)
	if err != nil {
		return nil, s.log.Error().LogErrorF("credentialID invalid - %w", err)
	}

	// _, err = uuid.Parse(session.TenantID)
	// if err != nil {
	// 	return nil, s.log.Error().LogErrorF("identityID invalid - %w", err)
	// }

	return &session, nil
}
