package zerotrust

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"gopkg.in/square/go-jose.v2"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey string

// SessionContextKey is the context key for the Login Session
const SessionContextKey contextKey = "session"

// GatewayMiddleware - Handles authenticating a request came from the authn services
type GatewayMiddleware struct {
	time       stime.TimeService
	publicKeys jose.JSONWebKeySet
}

// NewGatewayMiddleware - Generates a default AuthnMiddleware for use with authenticating a request came from the authn services
func NewGatewayMiddleware(time stime.TimeService, publicKeyLoader webkeys.WebKeysService) (*GatewayMiddleware, error) {
	publicKeys, err := publicKeyLoader.FetchJwks()
	if err != nil {
		return nil, err
	}

	return &GatewayMiddleware{
		time:       time,
		publicKeys: *publicKeys,
	}, nil
}

// Handler - Generates the handler you use to wrap the http routes
func (s *GatewayMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.FromRequest(r)
		if err != nil {
			fmt.Println("Session Middleware Failure", err)
			w.WriteHeader(404)
			return
		}

		// Don't really like using this map of any objects in the context for this, but it seems how its done.
		ctx := context.WithValue(r.Context(), SessionContextKey, &session.Session)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromRequest - Pulls out authenticationd details from the Request and calls Parse.
func (s *GatewayMiddleware) FromRequest(r *http.Request) (*SessionJwt, error) {
	authHeader, err := s.fromAuthHeader(r)
	if err != nil {
		return nil, err
	}

	session, err := s.Parse(authHeader)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *GatewayMiddleware) fromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header missing")
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// Parse - Parses the JWT token and verifies the signature came from AuthN via the public keys we obtain
func (s *GatewayMiddleware) Parse(tokenString string) (*SessionJwt, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SessionJwt{}, func(token *jwt.Token) (interface{}, error) {

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

	if claims, ok := token.Claims.(*SessionJwt); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Token is invalid")
}
