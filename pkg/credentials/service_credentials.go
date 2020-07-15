package credentials

import (
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// CredentialsService is a service that implents the logic for the CredentialsApiServicer
// This service should implement the business logic for every endpoint for the CredentialsApi API.
// Include any external packages or services that will be required by this service.
type CredentialsService struct {
	time       stime.TimeService
	repository CredentialRepository
}

// NewCredentialsService creates a default api service
func NewCredentialsService(time stime.TimeService, repository CredentialRepository) *CredentialsService {
	return &CredentialsService{
		time:       time,
		repository: repository,
	}
}

// DisableCredentials - Disables a credential so it can&#39;t be used anymore to login
func (s *CredentialsService) DisableCredentials(auth tmw.TumblerClaims, identityID string, credentialID string) (*api.Credential, error) {
	cred, err := s.repository.get(identityID, credentialID, auth.TenantID.String())
	if err != nil {
		return nil, err
	}

	caller := auth.IdentityID.String()
	now := s.time.Now()
	cred.DisabledOn = &now
	cred.DisabledBy = &caller

	saved, err := s.repository.update(*cred)
	if err != nil {
		return nil, err
	}

	// @TODO send notification to the email to notify them?

	return saved, nil
}

// ListCredentials - List the credentials this user has used.
func (s *CredentialsService) ListCredentials(auth tmw.TumblerClaims, identityID string) ([]api.Credential, error) {
	return s.repository.list(identityID, auth.TenantID.String())
}

func (s *CredentialsService) Login(login api.Login, nonce string, ip string) (*api.Credential, error) {
	// look into the repo for any matches
	cred, err := s.repository.lookup(login.CredentialID, login.TenantID)
	if err != nil {
		return nil, err
	}

	cred.LastUsedOn = s.time.Now()

	// Record the login happened and that the nonce is unique.
	err = s.repository.record(cred.CredentialID, nonce, ip, cred.LastUsedOn)
	if err != nil {
		return nil, err
	}

	saved, err := s.repository.update(*cred)
	if err != nil {
		return nil, err
	}

	// @TODO record login in a queue

	return saved, nil
}

func (s *CredentialsService) Register(identityID string, credentialID, tenantID string) (*api.Credential, error) {
	cred := api.Credential{
		CredentialID: credentialID,
		IdentityID:   identityID,
		TenantID:     tenantID,
		CreatedOn:    s.time.Now(),
		LastUsedOn:   s.time.Now(),
		DisabledBy:   nil,
		DisabledOn:   nil,
	}

	saved, err := s.repository.add(cred)
	if err != nil {
		return nil, err
	}

	// @TODO record new registered credential

	// @TODO email that a new credential was registered

	return saved, nil
}
