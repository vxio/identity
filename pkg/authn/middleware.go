package authn

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"gopkg.in/square/go-jose.v2"
)

// Middleware - Handles authenticating a request came from the authn services
type Middleware struct {
	time       stime.TimeService
	publicKeys jose.JSONWebKeySet
}

// NewMiddleware - Generates a default AuthnMiddleware for use with authenticating a request came from the authn services
func NewMiddleware(time stime.TimeService, publicKeyLoader webkeys.WebKeysService) (*Middleware, error) {
	publicKeys, err := publicKeyLoader.FetchJwks()
	if err != nil {
		return nil, err
	}

	return &Middleware{
		time:       time,
		publicKeys: *publicKeys,
	}, nil
}

// Handler - Generates the handler you use to wrap the http routes
func (s *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.FromRequest(r)
		if err != nil {
			return
		}

		// Don't really like using this map of any objects in the context for this, but it seems how its done.
		ctx := context.WithValue(r.Context(), loginSessionContextKey, session)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromRequest - Pulls out authenticationd details from the Request and calls Parse.
func (s *Middleware) FromRequest(r *http.Request) (*LoginSession, error) {
	cookie, err := r.Cookie("moov-authn")
	if err != nil {
		return nil, errors.New("No session found")
	}

	session, err := s.Parse(cookie.Value)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Parse - Parses the JWT token and verifies the signature came from AuthN via the public keys we obtain
func (s *Middleware) Parse(tokenString string) (*LoginSession, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LoginSession{}, func(token *jwt.Token) (interface{}, error) {

		// get the key ID `kid` from the jwt.Token
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not specified")
		}

		// search the returned keys from the JWKS
		k := s.publicKeys.Key(kid)
		if len(k) == 0 {
			return nil, errors.New("Could not find the kid in the public web key set")
		}

		return k[0].Key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*LoginSession); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Token is invalid")
}
