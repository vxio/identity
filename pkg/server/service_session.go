package identityserver

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/jwks"
	"gopkg.in/square/go-jose.v2"
)

/*


	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{

	})
*/

type TokenService interface {
	Generate(identityID string) (string, error)
}

type tokenService struct {
	time       TimeService
	jwks       jwks.JwksService
	expiration time.Duration
}

func NewTokenService(time TimeService, jwks jwks.JwksService, expiration time.Duration) TokenService {
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
		ExpiresAt: s.time.Now().Add(s.expiration).Unix(),
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

func getPrivateKey(keys *jose.JSONWebKeySet) *jose.JSONWebKey {
	for _, k := range keys.Keys {
		if !k.IsPublic() {
			return &k
		}
	}

	return nil
}
