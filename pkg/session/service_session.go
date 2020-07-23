package session

import (
	"net/http"
	"time"

	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/tumbler/pkg/jwe"
)

// SessionService - Generates the tokens for their fully logged in session.
type SessionService interface {
	Generate(r *http.Request, Session Session) (string, error)
	GenerateCookie(r *http.Request, session Session) (*http.Cookie, error)
}

type sessionService struct {
	time       stime.TimeService
	jweService jwe.JWEService
	expiration time.Duration
}

// NewSessionService - Creates a default instance of a SessionService
func NewSessionService(time stime.TimeService, jweService jwe.JWEService, config Config) SessionService {
	return &sessionService{
		time:       time,
		jweService: jweService,
		expiration: config.Expiration,
	}
}

// Generate - Creates the token string
func (s *sessionService) Generate(r *http.Request, session Session) (string, error) {
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
func (s *sessionService) GenerateCookie(r *http.Request, session Session) (*http.Cookie, error) {
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

func (s *sessionService) calculateExpiration() time.Time {
	return s.time.Now().Add(s.expiration)
}
