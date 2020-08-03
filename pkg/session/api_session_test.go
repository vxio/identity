package session_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/client"
)

func Test_SessionEndpoint(t *testing.T) {
	s := NewSessionScope(t)

	identity, err := s.identities.Register(client.Register{
		CredentialID: uuid.New().String(),
		TenantID:     s.claims.TenantID.String(),
		FirstName:    "John",
		LastName:     "Doe",
		NickName:     nil,
	}, nil)
	s.assert.Nil(err)

	s.claims.IdentityID = uuid.MustParse(identity.IdentityID)

	api := s.APIClient()

	_, resp, err := api.SessionApi.GetSessionDetails(context.Background())
	s.assert.Equal(200, resp.StatusCode)
	s.assert.Nil(err)
}

func Test_SessionEndpoint_Identity_Not_Found(t *testing.T) {
	s := NewSessionScope(t)

	api := s.APIClient()

	_, resp, err := api.SessionApi.GetSessionDetails(context.Background())
	s.assert.Equal(404, resp.StatusCode)
	s.assert.NotNil(err)
}
