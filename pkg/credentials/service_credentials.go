package credentials

import (
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// CredentialsApiServicer defines the api actions for the CredentialsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CredentialsService interface {
	DisableCredentials(tmw.TumblerClaims, string, string) (*client.Credential, error)
	ListCredentials(tmw.TumblerClaims, string) ([]client.Credential, error)

	Register(string, string, string) (*client.Credential, error)
	Login(client.Login, string, string) (*client.Credential, error)
}

// CredentialsService is a service that implents the logic for the CredentialsApiServicer
// This service should implement the business logic for every endpoint for the CredentialsApi API.
// Include any external packages or services that will be required by this service.
type credentialsService struct {
	time       stime.TimeService
	repository CredentialRepository
}

// NewCredentialsService creates a default api service
func NewCredentialsService(time stime.TimeService, repository CredentialRepository) CredentialsService {
	return &credentialsService{
		time:       time,
		repository: repository,
	}
}

// DisableCredentials - Disables a credential so it can&#39;t be used anymore to login
func (s *credentialsService) DisableCredentials(auth tmw.TumblerClaims, identityID string, credentialID string) (*client.Credential, error) {
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
func (s *credentialsService) ListCredentials(auth tmw.TumblerClaims, identityID string) ([]client.Credential, error) {
	return s.repository.list(identityID, auth.TenantID.String())
}

func (s *credentialsService) Login(login client.Login, nonce string, ip string) (*client.Credential, error) {

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

func (s *credentialsService) Register(identityID string, credentialID, tenantID string) (*client.Credential, error) {
	cred := client.Credential{
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
