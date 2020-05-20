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

type AuthnMiddleware struct {
	time       stime.TimeService
	publicKeys jose.JSONWebKeySet
}

func NewAuthnMiddleware(time stime.TimeService, publicKeyLoader webkeys.WebKeysService) (*AuthnMiddleware, error) {
	publicKeys, err := publicKeyLoader.FetchJwks()
	if err != nil {
		return nil, err
	}

	return &AuthnMiddleware{
		time:       time,
		publicKeys: *publicKeys,
	}, nil
}

func (s *AuthnMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.FromRequest(r)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), "LoginSession", session)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *AuthnMiddleware) FromRequest(r *http.Request) (*LoginSession, error) {
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

func (s *AuthnMiddleware) Parse(tokenString string) (*LoginSession, error) {
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
	} else {
		return nil, errors.New("Token is invalid")
	}
}
