package credentials_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func Test_DisableAPI(t *testing.T) {
	a, s := Setup(t)

	cred, err := s.RegisterRandom()
	a.Nil(err)

	resp, err := s.api.CredentialsApi.DisableCredentials(context.Background(), cred.IdentityID, cred.CredentialID)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)
}

func Test_DisableAPI_NotFound(t *testing.T) {
	a, s := Setup(t)

	cred, err := s.RegisterRandom()
	a.Nil(err)

	resp, err := s.api.CredentialsApi.DisableCredentials(context.Background(), cred.IdentityID, uuid.New().String())
	a.NotNil(err)
	a.Equal(404, resp.StatusCode)

	resp, err = s.api.CredentialsApi.DisableCredentials(context.Background(), uuid.New().String(), cred.CredentialID)
	a.NotNil(err)
	a.Equal(404, resp.StatusCode)
}

func Test_ListAPI(t *testing.T) {
	a, s := Setup(t)

	cred, err := s.RegisterRandom()
	a.Nil(err)

	found, resp, err := s.api.CredentialsApi.ListCredentials(context.Background(), cred.IdentityID)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
	a.Len(found, 1)

	a.Contains(found, *cred)
}

func Test_ListAPI_NotFound(t *testing.T) {
	a, s := Setup(t)

	_, err := s.RegisterRandom()
	a.Nil(err)

	found, resp, err := s.api.CredentialsApi.ListCredentials(context.Background(), uuid.New().String())
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
	a.Len(found, 0)
}
