package identities

import (
	"errors"

	"github.com/google/uuid"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// Service - Service that implents the logic for the IdentitiesApiServicer
// This service should implement the business logic for every endpoint for the IdentitiesApi API.
// Include any external packages or services that will be required by this service.
type Service interface {
	DisableIdentity(claims tmw.TumblerClaims, identityID string) error
	GetIdentity(claims tmw.TumblerClaims, identityID string) (*client.Identity, error)
	ListIdentities(claims tmw.TumblerClaims, orgID string) ([]client.Identity, error)
	UpdateIdentity(claims tmw.TumblerClaims, identityID string, update client.UpdateIdentity) (*client.Identity, error)

	Register(register client.Register, invite *client.Invite) (*client.Identity, error)
	GetIdentityByID(identityID string) (*client.Identity, error)
}

type service struct {
	time       stime.TimeService
	repository Repository
}

// NewIdentitiesService creates a default service
func NewIdentitiesService(time stime.TimeService, repository Repository) Service {
	return &service{
		time:       time,
		repository: repository,
	}
}

// DisableIdentity - Disable an identity. Its left around for historical reporting
func (s *service) DisableIdentity(claims tmw.TumblerClaims, identityID string) error {
	identity, err := s.GetIdentity(claims, identityID)
	if err != nil {
		return err
	}

	now := s.time.Now()
	callerIdentityID := claims.IdentityID.String()

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
func (s *service) GetIdentity(claims tmw.TumblerClaims, identityID string) (*client.Identity, error) {
	i, e := s.GetIdentityByID(identityID)
	if e != nil {
		return nil, e
	}

	if i.TenantID != claims.TenantID.String() {
		return nil, errors.New("TenantID of session user doesn't match retrieved identity")
	}

	return i, nil
}

// ListIdentities - List identities and associates userId
func (s *service) ListIdentities(claims tmw.TumblerClaims, orgID string) ([]client.Identity, error) {
	identities, err := s.repository.list(api.TenantID(claims.TenantID))
	return identities, err
}

// UpdateIdentity - Update a specific Identity
func (s *service) UpdateIdentity(claims tmw.TumblerClaims, identityID string, update client.UpdateIdentity) (*client.Identity, error) {
	identity, err := s.GetIdentity(claims, identityID)
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

	identity.Phones = []client.Phone{}
	for _, p := range update.Phones {
		_, err := uuid.Parse(p.PhoneID)
		if err != nil {
			p.PhoneID = uuid.New().String()
		}

		identity.Phones = append(
			identity.Phones,
			client.Phone{
				IdentityID: identity.IdentityID,
				PhoneID:    p.PhoneID,
				Number:     p.Number,
				Validated:  p.Validated,
				Type:       p.Type,
			},
		)
	}

	identity.Addresses = []client.Address{}
	for _, a := range update.Addresses {
		_, err := uuid.Parse(a.AddressID)
		if err != nil {
			a.AddressID = uuid.New().String()
		}

		identity.Addresses = append(
			identity.Addresses,
			client.Address{
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
func (s *service) Register(register client.Register, invite *client.Invite) (*client.Identity, error) {
	identityID := uuid.New().String()

	phones := []client.Phone{}
	for _, rp := range register.Phones {
		phones = append(phones, client.Phone{
			IdentityID: identityID,
			PhoneID:    uuid.New().String(),
			Number:     rp.Number,
			Validated:  false,
			Type:       rp.Type,
		})
	}

	addresses := []client.Address{}
	for _, ra := range register.Addresses {
		addresses = append(addresses, client.Address{
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

	identity := client.Identity{
		IdentityID:    identityID,
		TenantID:      register.TenantID,
		InviteID:      nil,
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
		LastLogin:     client.LastLogin{},
		LastUpdatedOn: s.time.Now(),
	}

	if invite != nil {
		identity.TenantID = invite.TenantID
		identity.InviteID = &invite.InviteID
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
func (s *service) GetIdentityByID(identityID string) (*client.Identity, error) {
	i, e := s.repository.get(identityID)
	if e != nil {
		return nil, e
	}

	return i, nil
}
