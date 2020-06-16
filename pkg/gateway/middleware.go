package gateway

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"gopkg.in/square/go-jose.v2"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey string

// SessionContextKey is the context key for the Login Session
const SessionContextKey contextKey = "session"

// Middleware - Handles authenticating a request came from the authn services
type Middleware struct {
	log        logging.Logger
	time       stime.TimeService
	publicKeys jose.JSONWebKeySet
}

// NewMiddleware - Generates a default AuthnMiddleware for use with authenticating a request came from the authn services
func NewMiddleware(log logging.Logger, time stime.TimeService, publicKeyLoader webkeys.WebKeysService) (*Middleware, error) {
	publicKeys, err := publicKeyLoader.Keys()
	if err != nil {
		return nil, err
	}

	return &Middleware{
		log:        log,
		time:       time,
		publicKeys: *publicKeys,
	}, nil
}

// Handler - Generates the handler you use to wrap the http routes
func (s *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.FromRequest(r)
		if err != nil {
			s.log.Error().LogError("Gateway Token Failure", err)
			w.WriteHeader(404)
			return
		}

		// Don't really like using this map of any objects in the context for this, but it seems how its done.
		ctx := context.WithValue(r.Context(), SessionContextKey, session)

		h.ServeHTTP(w, r.Clone(ctx))
	})
}

// FromRequest - Pulls out authenticationd details from the Request and calls Parse.
func (s *Middleware) FromRequest(r *http.Request) (*Session, error) {
	authHeader, err := s.fromAuthHeader(r)
	if err != nil {
		return nil, err
	}

	claims, err := s.Parse(authHeader)
	if err != nil {
		return nil, err
	}

	session := Session{
		CallerID: IdentityID(claims.CallerID),
		TenantID: TenantID(claims.TenantID),
	}

	return &session, nil
}

func (s *Middleware) fromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("X-Auth-Token")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// Parse - Parses the JWT token and verifies the signature came from AuthN via the public keys we obtain
func (s *Middleware) Parse(tokenString string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {

		// get the key ID `kid` from the jwt.Token
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not specified")
		}

		// search the returned keys from the JWKS
		k := s.publicKeys.Key(kid)
		if len(k) == 0 {
			return nil, errors.New("could not find the kid in the public web key set")
		}

		return k[0].Key, nil
	})

	if err != nil {
		return nil, s.log.Error().LogErrorF("unable to parse gateway token - %w", err)
	}

	claims, ok := token.Claims.(*SessionClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
