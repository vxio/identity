package session

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/logging"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// SessionService - Generates the tokens for their fully logged in session.
type SessionService interface {
	GetDetails(claims tmw.TumblerClaims) (*client.SessionDetails, error)
	ChangeDetails(req *http.Request, claims tmw.TumblerClaims, updates client.ChangeSessionDetails) (*client.SessionDetails, *http.Cookie, error)
}

type sessionService struct {
	logger      logging.Logger
	identities  identities.Service
	service     TokenService
	credentials credentials.CredentialsService
	config      Config
}

// NewSessionService - Creates a default instance of a SessionService
func NewSessionService(logger logging.Logger, identities identities.Service, service TokenService, credentials credentials.CredentialsService, config Config) SessionService {
	return &sessionService{
		logger:      logger,
		identities:  identities,
		service:     service,
		credentials: credentials,
		config:      config,
	}
}

func (s *sessionService) GetDetails(claims tmw.TumblerClaims) (*client.SessionDetails, error) {
	// Require people have used credentials to login
	if err := s.ensureHuman(claims); err != nil {
		return nil, err
	}

	identity, err := s.getIdentity(claims.IdentityID)
	if err != nil {
		return nil, s.logger.Info().LogError("Unable to lookup identity", err)
	}

	details := client.SessionDetails{
		CredentialID: claims.CredentialID.String(),
		TenantID:     claims.TenantID.String(),
		IdentityID:   claims.IdentityID.String(),
		ExpiresIn:    claims.Expiry.Time().Unix(),
		FirstName:    identity.FirstName,
		LastName:     identity.LastName,
		NickName:     identity.NickName,
		ImageUrl:     identity.ImageUrl,
	}

	return &details, nil
}

func (s *sessionService) ChangeDetails(req *http.Request, claims tmw.TumblerClaims, update client.ChangeSessionDetails) (*client.SessionDetails, *http.Cookie, error) {
	if !s.config.EnablePutSession {
		return nil, nil, ErrPutSessionNotEnabled
	}

	// Require people have used credentials to login
	if err := s.ensureHuman(claims); err != nil {
		return nil, nil, err
	}

	// Change the session details here.
	if update.TenantID != nil {
		tid, err := uuid.Parse(*update.TenantID)
		if err != nil {
			return nil, nil, err
		}

		claims.TenantID = tid
	}

	// Get the new details after the changes
	details, err := s.GetDetails(claims)
	if err != nil {
		return nil, nil, err
	}

	// Record the use of the credential to log into another tenant/identity
	if err := s.credentials.Record(claims.CredentialID.String(), claims.TenantID.String(), claims.ID, claims.RemoteAddr); err != nil {
		return nil, nil, s.logger.LogError("Already logged into this tenant, failing", err)
	}

	// Generate a new session
	cookie, err := s.service.GenerateCookie(req, Session{
		IdentityID:   *claims.IdentityID,
		TenantID:     claims.TenantID,
		CredentialID: *claims.CredentialID,
	})
	if err != nil {
		return nil, nil, s.logger.Error().LogError("Unable to generate cookie", err)
	}

	return details, cookie, nil
}

func (s *sessionService) getIdentity(identityID *uuid.UUID) (*client.Identity, error) {
	if identityID == nil {
		return nil, ErrIdentityNotFound
	}

	identity, err := s.identities.GetIdentityByID(identityID.String())
	if err != nil {
		return nil, s.logger.Info().LogError("Unable to lookup identity", err)
	}

	return identity, nil
}

func (s *sessionService) ensureHuman(claims tmw.TumblerClaims) error {
	if claims.CredentialID == nil {
		return ErrCredentialsNotSet
	}

	if claims.IdentityID == nil {
		return ErrIdentityNotFound
	}

	return nil
}
