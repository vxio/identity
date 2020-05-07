/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package identities

import (
	"errors"

	"github.com/google/uuid"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/utils"
)

// IdentitiesApiService is a service that implents the logic for the IdentitiesApiServicer
// This service should implement the business logic for every endpoint for the IdentitiesApi API.
// Include any external packages or services that will be required by this service.
type IdentitiesService struct {
	time       utils.TimeService
	repository IdentityRepository
}

// NewIdentitiesApiService creates a default api service
func NewIdentitiesService(time utils.TimeService, repository IdentityRepository) *IdentitiesService {
	return &IdentitiesService{
		time:       time,
		repository: repository,
	}
}

// DisableIdentity - Disable an identity. Its left around for historical reporting
func (s *IdentitiesService) DisableIdentity(session api.Session, identityID string) error {
	identity, err := s.repository.get(session.TenantID, identityID)
	if err != nil {
		return err
	}

	now := s.time.Now()
	callerIdentityID := session.CallerID.String()

	identity.DisabledOn = &now
	identity.DisabledBy = &callerIdentityID
	identity.LastUpdatedOn = s.time.Now()

	_, nil := s.repository.update(*identity)
	if err != nil {
		return err
	}

	// supposed to be 204 no content...
	return nil
}

// GetIdentity - List identities and associates userId
func (s *IdentitiesService) GetIdentity(session api.Session, identityID string) (*api.Identity, error) {
	i, e := s.repository.get(session.TenantID, identityID)
	if e != nil {
		return nil, errors.New("Identity not found")
	}

	return i, nil
}

// ListIdentities - List identities and associates userId
func (s *IdentitiesService) ListIdentities(session api.Session, orgID string) ([]api.Identity, error) {
	identities, err := s.repository.list(session.TenantID)
	return identities, err
}

// UpdateIdentity - Update a specific Identity
func (s *IdentitiesService) UpdateIdentity(session api.Session, identityID string, update api.UpdateIdentity) (*api.Identity, error) {
	identity, err := s.repository.get(session.TenantID, identityID)
	if err != nil {
		return nil, err
	}

	identity.FirstName = update.FirstName
	identity.MiddleName = update.MiddleName
	identity.LastName = update.LastName
	identity.NickName = update.NickName
	identity.Suffix = update.Suffix
	identity.BirthDate = update.BirthDate
	identity.Status = update.Status
	identity.Phones = update.Phones
	identity.Addresses = update.Addresses
	identity.LastUpdatedOn = s.time.Now()

	updated, err := s.repository.update(*identity)
	if err != nil {
		return nil, err
	}

	// @TODO record update and email identity that changes were made.

	return updated, err
}

func (s *IdentitiesService) Register(register api.Register) (*api.Identity, error) {
	identity := api.Identity{
		IdentityID:    uuid.New().String(),
		TenantID:      register.TenantID,
		FirstName:     register.FirstName,
		MiddleName:    register.MiddleName,
		LastName:      register.LastName,
		NickName:      register.NickName,
		Suffix:        register.Suffix,
		BirthDate:     register.BirthDate,
		Status:        "none",
		Email:         register.Email,
		EmailVerified: false,
		Phones:        register.Phones,
		Addresses:     register.Addresses,
		RegisteredOn:  s.time.Now(),
		LastLogin:     api.LastLogin{},
		LastUpdatedOn: s.time.Now(),
	}

	saved, err := s.repository.add(identity)
	if err != nil {
		return nil, err
	}

	// @TODO record user was registered

	// @TODO send email notification to get registered email verified

	return saved, nil
}
