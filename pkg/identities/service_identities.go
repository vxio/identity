package identities

import (
	"errors"

	"github.com/google/uuid"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/stime"
)

// Service - Service that implents the logic for the IdentitiesApiServicer
// This service should implement the business logic for every endpoint for the IdentitiesApi API.
// Include any external packages or services that will be required by this service.
type Service struct {
	time       stime.TimeService
	repository Repository
}

// NewIdentitiesService creates a default service
func NewIdentitiesService(time stime.TimeService, repository Repository) *Service {
	return &Service{
		time:       time,
		repository: repository,
	}
}

// DisableIdentity - Disable an identity. Its left around for historical reporting
func (s *Service) DisableIdentity(session gateway.Session, identityID string) error {
	identity, err := s.GetIdentity(session, identityID)
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
func (s *Service) GetIdentity(session gateway.Session, identityID string) (*api.Identity, error) {
	i, e := s.GetIdentityByID(identityID)
	if e != nil {
		return nil, e
	}

	if i.TenantID != session.TenantID.String() {
		return nil, errors.New("TenantID of session user doesn't match retrieved identity")
	}

	return i, nil
}

// ListIdentities - List identities and associates userId
func (s *Service) ListIdentities(session gateway.Session, orgID string) ([]api.Identity, error) {
	identities, err := s.repository.list(session.TenantID)
	return identities, err
}

// UpdateIdentity - Update a specific Identity
func (s *Service) UpdateIdentity(session gateway.Session, identityID string, update api.UpdateIdentity) (*api.Identity, error) {
	identity, err := s.GetIdentity(session, identityID)
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
	identity.LastUpdatedOn = s.time.Now()

	identity.Phones = []api.Phone{}
	for _, p := range update.Phones {
		_, err := uuid.Parse(p.PhoneID)
		if err != nil {
			p.PhoneID = uuid.New().String()
		}

		identity.Phones = append(
			identity.Phones,
			api.Phone{
				IdentityID: identity.IdentityID,
				PhoneID:    p.PhoneID,
				Number:     p.Number,
				Validated:  p.Validated,
				Type:       p.Type,
			},
		)
	}

	identity.Addresses = []api.Address{}
	for _, a := range update.Addresses {
		_, err := uuid.Parse(a.AddressID)
		if err != nil {
			a.AddressID = uuid.New().String()
		}

		identity.Addresses = append(
			identity.Addresses,
			api.Address{
				IdentityID: identity.IdentityID,
				AddressID:  a.AddressID,
				Type:       a.Type,
				Address1:   a.Address1,
				Address2:   a.Address2,
				City:       a.City,
				State:      a.State,
				PostalCode: a.PostalCode,
				Country:    a.Country,
				Validated:  a.Validated,
			},
		)
	}

	updated, err := s.repository.update(*identity)
	if err != nil {
		return nil, err
	}

	// @TODO record update and email identity that changes were made.

	return updated, err
}

// Register - Takes an invite and the registration information and creates the new identity from it.
func (s *Service) Register(invite api.Invite, register api.Register) (*api.Identity, error) {
	identityID := uuid.New().String()

	phones := []api.Phone{}
	for _, rp := range register.Phones {
		phones = append(phones, api.Phone{
			IdentityID: identityID,
			PhoneID:    uuid.New().String(),
			Number:     rp.Number,
			Validated:  false,
			Type:       rp.Type,
		})
	}

	addresses := []api.Address{}
	for _, ra := range register.Addresses {
		addresses = append(addresses, api.Address{
			IdentityID: identityID,
			AddressID:  uuid.New().String(),
			Type:       ra.Type,
			Address1:   ra.Address1,
			Address2:   ra.Address2,
			City:       ra.City,
			State:      ra.State,
			PostalCode: ra.PostalCode,
			Country:    ra.Country,
			Validated:  false,
		})
	}

	identity := api.Identity{
		IdentityID:    identityID,
		TenantID:      invite.TenantID,
		InviteID:      invite.InviteID,
		FirstName:     register.FirstName,
		MiddleName:    register.MiddleName,
		LastName:      register.LastName,
		NickName:      register.NickName,
		Suffix:        register.Suffix,
		BirthDate:     register.BirthDate,
		Status:        "none",
		Email:         register.Email,
		EmailVerified: false,
		Phones:        phones,
		Addresses:     addresses,
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

// GetIdentityByID - Returns the Identity specified by the ID. Used after a login session to get identity information
func (s *Service) GetIdentityByID(identityID string) (*api.Identity, error) {
	i, e := s.repository.get(identityID)
	if e != nil {
		return nil, e
	}

	return i, nil
}
