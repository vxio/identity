package session

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

// SessionService - Generates the tokens for their fully logged in session.
type SessionService interface {
	Generate(Session Session) (string, error)
	GenerateCookie(session Session) (*http.Cookie, error)

	FromRequest(r *http.Request) (*Session, error)
	Parse(tokenString string) (*Session, error)
}

type sessionService struct {
	time       stime.TimeService
	keys       webkeys.WebKeysService
	expiration time.Duration
}

// NewSessionService - Creates a default instance of a SessionService
func NewSessionService(time stime.TimeService, keys webkeys.WebKeysService, config Config) SessionService {
	return &sessionService{
		time:       time,
		keys:       keys,
		expiration: config.Expiration,
	}
}

// Generate - Creates the token string
func (s *sessionService) Generate(session Session) (string, error) {
	keys, err := s.keys.Keys()
	if err != nil {
		return "", err
	}

	privateKey := getPrivateKey(keys)
	if privateKey == nil {
		return "", errors.New("Unable to find a private key to use")
	}

	signingMethod := jwt.GetSigningMethod(privateKey.Algorithm)

	sessionJwt := SessionJwt{
		Session: session,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: s.calculateExpiration().Unix(),
			NotBefore: s.time.Now().Add(time.Minute * -1).Unix(),
			IssuedAt:  s.time.Now().Unix(),
			Id:        uuid.New().String(),
			Subject:   session.IdentityID.String(),

			Audience: "moov",
			Issuer:   "moov",
		},
	}

	token := jwt.NewWithClaims(signingMethod, sessionJwt)
	token.Header["kid"] = privateKey.KeyID

	tokenString, err := token.SignedString(privateKey.Key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateCookie - Generates the token and the cookie version of it.
func (s *sessionService) GenerateCookie(session Session) (*http.Cookie, error) {
	value, err := s.Generate(session)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     "moov",
		Value:    value,
		Path:     "/",
		Expires:  s.calculateExpiration(),
		MaxAge:   int(s.expiration.Seconds()),
		SameSite: http.SameSiteDefaultMode,
		Secure:   false,
		HttpOnly: true,
	}, nil
}

// FromRequest - Pulls out authenticationd details from the Request and calls Parse.
func (s *sessionService) FromRequest(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("moov")
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
func (s *sessionService) Parse(tokenString string) (*Session, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SessionJwt{}, func(token *jwt.Token) (interface{}, error) {

		// get the key ID `kid` from the jwt.Token
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not specified")
		}

		keys, err := s.keys.Keys()
		if err != nil {
			return nil, err
		}

		// search the returned keys from the JWKS
		found := keys.Key(kid)

		for _, k := range found {
			if k.IsPublic() {
				return k.Key, nil
			}
		}

		return nil, errors.New("Could not find the kid in the public web key set")
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SessionJwt); ok && token.Valid {
		return &claims.Session, nil
	}

	return nil, errors.New("Token is invalid")
}

func (s *sessionService) calculateExpiration() time.Time {
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

func getPublicKey(keys *jose.JSONWebKeySet) *jose.JSONWebKey {
	for _, k := range keys.Keys {
		if k.IsPublic() {
			return &k
		}
	}

	return nil
}
