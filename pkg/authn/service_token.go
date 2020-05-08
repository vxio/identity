package authn

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"gopkg.in/square/go-jose.v2"
)

type TokenService interface {
	Generate(identityID string) (string, error)
	GenerateCookie(identityID string) (*http.Cookie, error)
}

type tokenService struct {
	time       stime.TimeService
	jwks       webkeys.WebKeysService
	expiration time.Duration
}

func NewTokenService(time stime.TimeService, jwks webkeys.WebKeysService, expiration time.Duration) TokenService {
	return &tokenService{
		time:       time,
		jwks:       jwks,
		expiration: expiration,
	}
}

// DeleteInvite - Delete an invite that was sent and invalidate the token.
func (s *tokenService) Generate(identityID string) (string, error) {
	keys, err := s.jwks.FetchJwks()
	if err != nil {
		return "", err
	}

	privateKey := getPrivateKey(keys)
	if privateKey == nil {
		return "", errors.New("Unable to find a private key to use")
	}

	signingMethod := jwt.GetSigningMethod(privateKey.Algorithm)

	//jwt.SigningMethodRS256
	token := jwt.NewWithClaims(signingMethod, jwt.StandardClaims{
		ExpiresAt: s.calculateExpiration().Unix(),
		NotBefore: s.time.Now().Add(time.Minute * -1).Unix(),
		IssuedAt:  s.time.Now().Unix(),
		Id:        uuid.New().String(),
		Subject:   identityID,

		Audience: "moov",
		Issuer:   "moov",
	})

	tokenString, err := token.SignedString(privateKey.Key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *tokenService) GenerateCookie(identityID string) (*http.Cookie, error) {
	value, err := s.Generate(identityID)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     "moov",
		Value:    value,
		Path:     "/",
		Expires:  s.calculateExpiration(),
		MaxAge:   int(s.expiration.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
		HttpOnly: true,
	}, nil
}

func (s *tokenService) calculateExpiration() time.Time {
	return s.time.Now().Add(s.expiration)
}

func getPrivateKey(keys *jose.JSONWebKeySet) *jose.JSONWebKey {
	for _, k := range keys.Keys {
		if !k.IsPublic() {
			return &k
		}
	}

	return nil
}
