package session_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/client"
)

func Test_SessionEndpoint(t *testing.T) {
	s := NewSessionScope(t)
	imageUrl := "http://kittens.com/123.jpg"
	identity, err := s.identities.Register(client.Register{
		CredentialID: uuid.New().String(),
		TenantID:     s.claims.TenantID.String(),
		FirstName:    "John",
		LastName:     "Doe",
		NickName:     nil,
		ImageUrl:     &imageUrl,
	})
	s.assert.Nil(err)

	iid := uuid.MustParse(identity.IdentityID)
	s.claims.Subject = iid.String()
	s.claims.IdentityID = &iid

	api := s.APIClient()

	sessionDetails, resp, err := api.SessionApi.GetSessionDetails(context.Background())
	s.assert.Equal(200, resp.StatusCode)
	s.assert.Nil(err)
	s.assert.Equal(identity.ImageUrl, sessionDetails.ImageUrl)
}

func Test_SessionEndpoint_Identity_Not_Found(t *testing.T) {
	s := NewSessionScope(t)

	api := s.APIClient()

	_, resp, err := api.SessionApi.GetSessionDetails(context.Background())
	s.assert.Equal(404, resp.StatusCode)
	s.assert.NotNil(err)
}
