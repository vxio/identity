package authn

import (
	"errors"
	"net/http"

	identity "github.com/moov-io/identity/pkg/client"
	log "github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/tumbler/pkg/jwe"
)

// LoginSession is the values of the JWT coming in from the Authentication services.
type LoginSession struct {
	jwe.Claims

	// CSRF state token used during login
	State string `json:"st"`

	// Domain this was created under and only usable under.
	Origin string `json:"or"`

	// Flow this session was stated with and must end with
	Flow string `json:"fl"`

	// List of available providers for the tenantID
	Providers []string `json:"ps,omitempty"`

	// Provider that supplied the SubjectID
	Provider string `json:"pv,omitempty"`

	// Unique ID of the user under the external provider.
	SubjectID string `json:"si,omitempty"`

	// Set during logging in everytime and used to look up credentials
	Issuer *string `json:"pi"` // Issuer attribute of the login

	// IP Address of the login
	IP string `json:"ip"`

	// Scope of what this token is allow to do.
	Scopes []string `json:"scp"`

	// Login URL for the start of the flow
	LoginURL string `json:"lu"`

	// Store whatever we can get from the OIDC provider if the invite code isn't empty
	identity.Register
}

// LoginSessionFromRequest - Pulls the Login Session out of the context of a request
func LoginSessionFromRequest(r *http.Request) (*LoginSession, error) {
	session, ok := r.Context().Value(LoginSessionContextKey).(*LoginSession)
	if !ok || session == nil {
		return nil, errors.New("unable to find LoginSession in context")
	}
	return session, nil
}

// WithLoginSessionFromRequest - Pulls the Login Session out of the context of a request if its not available returns an error response on `w`.
func WithLoginSessionFromRequest(l log.Logger, w http.ResponseWriter, r *http.Request, scopes []string, run func(LoginSession)) {
	session, err := LoginSessionFromRequest(r)
	if err != nil {
		l.Error().LogError("LoginSessionFromRequest errored", err)
		w.WriteHeader(404)
		return
	}

	hasAllScopes := true
	for _, v := range scopes {
		found := false
		for _, s := range session.Scopes {
			if v == s {
				found = true
				break
			}
		}
		if !found {
			hasAllScopes = false
			break
		}
	}

	if !hasAllScopes {
		l.Error().LogError("LoginSessionFromRequest missing scopes", err)
		w.WriteHeader(404)
		return
	}

	run(*session)
}
