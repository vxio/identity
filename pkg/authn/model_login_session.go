package authn

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	api "github.com/moov-io/identity/pkg/api"
)

// LoginSession is the values of the JWT coming in from the Authentication services.
type LoginSession struct {
	State string `json:"state"` // CSRF state token used during login

	// Set during logging in everytime and used to look up credentials
	Issuer *string `json:"issuer"` // Issuer attribute of the login

	// IP Address of the login
	IP string `json:"ip"`

	// Scope of what this token is allow to do.
	Scopes []string `'json:"scp"`

	// standard JWT claims like expirations etc...
	jwt.StandardClaims

	// Store whatever we can get from the OIDC provider if the invite code isn't empty
	api.Register
}

func LoginSessionFromRequest(r *http.Request) (*LoginSession, error) {
	session, ok := r.Context().Value(LoginSessionContextKey).(*LoginSession)
	if !ok || session == nil {
		return nil, errors.New("Unable to find LoginSession in context")
	}
	return session, nil
}

func WithLoginSessionFromRequest(l log.Logger, w http.ResponseWriter, r *http.Request, run func(LoginSession)) {
	session, err := LoginSessionFromRequest(r)
	if err != nil {
		l.Log("level", "error", "msg", err.Error())
		w.WriteHeader(404)
		return
	}

	run(*session)
}
