package identityserver

import (
	"errors"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/moov-io/identity/pkg/jwks"
	"gopkg.in/square/go-jose.v2"
)

func NewJWTMiddleware(jwksLoader jwks.JwksService) (*jwtmiddleware.JWTMiddleware, error) {

	// Fetch the JWKS from our source.
	jwks, err := jwksLoader.FetchJwks()
	if err != nil {
		return nil, errors.New("Unable to load the jwks")
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			// get the key ID `kid` from the jwt.Token
			kid := token.Header["kid"].(string)

			// search the returned keys from the JWKS
			k := jwks.Key(kid)
			if len(k) == 0 {
				return nil, errors.New("Could not find the kid in the jwks web key set")
			}

			return k[0], nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return jwtMiddleware, nil
}

func findKey(kid string, keys []jose.JSONWebKey) *jose.JSONWebKey {
	for _, k := range keys {
		if k.KeyID == kid {
			return &k
		}
	}

	return nil
}
