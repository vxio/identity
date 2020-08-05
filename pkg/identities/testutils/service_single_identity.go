package identitiestestutils

import (
	"errors"

	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/identities"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

var ErrNotImplemented = errors.New("not implemented")

func NewSingleService(identity *client.Identity) identities.Service {
	if identity == nil {
		f := NewFuzzer()
		identity = &client.Identity{}
		f.Fuzz(identity)
	}

	return &singleService{identity: *identity}
}

type singleService struct {
	identity client.Identity
}

func (s *singleService) DisableIdentity(claims tmw.TumblerClaims, identityID string) error {
	panic(ErrNotImplemented)
}

func (s *singleService) GetIdentity(claims tmw.TumblerClaims, identityID string) (*client.Identity, error) {
	shallowCopy := s.identity
	shallowCopy.IdentityID = identityID
	return &shallowCopy, nil
}

func (s *singleService) ListIdentities(claims tmw.TumblerClaims, orgID string) ([]client.Identity, error) {
	return []client.Identity{s.identity}, nil
}

func (s *singleService) UpdateIdentity(claims tmw.TumblerClaims, identityID string, update client.UpdateIdentity) (*client.Identity, error) {
	panic(ErrNotImplemented)
}

func (s *singleService) Register(register client.Register, invite *client.Invite) (*client.Identity, error) {
	panic(ErrNotImplemented)
}

func (s *singleService) GetIdentityByID(identityID string) (*client.Identity, error) {
	shallowCopy := s.identity
	shallowCopy.IdentityID = identityID
	return &shallowCopy, nil
}
