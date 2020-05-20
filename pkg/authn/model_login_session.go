package authn

import (
	"github.com/dgrijalva/jwt-go"
	api "github.com/moov-io/identity/pkg/api"
)

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
