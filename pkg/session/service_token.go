package session

import (
	"net/http"
	"time"

	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/tumbler/pkg/jwe"
)

// TokenService - Generates the tokens for their fully logged in session.
type TokenService interface {
	Generate(r *http.Request, Session Session) (string, error)
	GenerateCookie(r *http.Request, session Session) (*http.Cookie, error)
}

type tokenService struct {
	time       stime.TimeService
	jweService jwe.JWEService
	expiration time.Duration
}

// NewTokenService - Creates a default instance of a SessionService
func NewTokenService(time stime.TimeService, jweService jwe.JWEService, config Config) TokenService {
	return &tokenService{
		time:       time,
		jweService: jweService,
		expiration: config.Expiration,
	}
}

// Generate - Creates the token string
func (s *tokenService) Generate(r *http.Request, session Session) (string, error) {
	c, err := s.jweService.Start(r)
	if err != nil {
		return "", err
	}

	tokenString, err := s.jweService.Serialize(c, session)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateCookie - Generates the token and the cookie version of it.
func (s *tokenService) GenerateCookie(r *http.Request, session Session) (*http.Cookie, error) {
	value, err := s.Generate(r, session)
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

func (s *tokenService) calculateExpiration() time.Time {
	return s.time.Now().Add(s.expiration)
}
