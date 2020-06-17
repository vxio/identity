package credentials_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	api "github.com/moov-io/identity/pkg/api"
)

func Test_Register(t *testing.T) {
	a, s := Setup(t)

	identityID := uuid.New().String()
	provider := "moovtest"
	subjectID := uuid.New().String()
	tenantID := uuid.New().String()

	cred, err := s.service.Register(identityID, provider, subjectID, tenantID)
	a.Nil(err)

	a.Equal(identityID, cred.IdentityID)
	a.Equal(provider, cred.Provider)
	a.Equal(subjectID, cred.SubjectID)

	a.Equal(s.time.Now(), cred.CreatedOn)
	a.Equal(s.time.Now(), cred.LastUsedOn)

	a.Nil(cred.DisabledBy)
	a.Nil(cred.DisabledOn)

	// register again should fail.
	_, err = s.service.Register(identityID, provider, subjectID, tenantID)
	a.NotNil(err)
}

func Test_List(t *testing.T) {
	a, s := Setup(t)

	cred, err := s.RegisterRandom()
	a.Nil(err)

	// Add noise
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()

	found, err := s.service.ListCredentials(s.session, cred.IdentityID)
	a.Nil(err)
	a.Len(found, 1)
	a.Contains(found, *cred)
}

func Test_Disable(t *testing.T) {
	a, s := Setup(t)

	cred, err := s.RegisterRandom()
	a.Nil(err)

	// Add noise
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	s.time.Change(s.time.Now().Add(time.Hour))

	disabled, err := s.service.DisableCredentials(s.session, cred.IdentityID, cred.CredentialID)
	a.Nil(err)

	a.Equal(cred.CredentialID, disabled.CredentialID)
	a.Equal(cred.IdentityID, disabled.IdentityID)
	a.Equal(cred.Provider, disabled.Provider)
	a.Equal(cred.SubjectID, disabled.SubjectID)
	a.Equal(cred.CreatedOn, disabled.CreatedOn)

	a.NotNil(disabled.DisabledBy)
	a.NotNil(disabled.DisabledOn)
	a.Equal(s.session.IdentityID.String(), *disabled.DisabledBy)
	a.Equal(s.time.Now(), *disabled.DisabledOn)

	found, err := s.service.ListCredentials(s.session, cred.IdentityID)
	a.Nil(err)
	a.Len(found, 1)
	a.Contains(found, *disabled)
}

func Test_Login(t *testing.T) {
	a, s := Setup(t)

	cred, err := s.RegisterRandom()
	a.Nil(err)

	// Add noise
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()

	login := api.Login{Provider: cred.Provider, SubjectID: cred.SubjectID, TenantID: cred.TenantID}
	nonce := uuid.New().String()
	ip := "1.2.3.4"

	// first login should work
	loggedIn, err := s.service.Login(login, nonce, ip)
	a.Nil(err)

	a.Equal(cred, loggedIn)

	// Can't login twice with the same nonce
	_, err = s.service.Login(login, nonce, ip)
	a.NotNil(err)

	// Difference nonce and works again
	nonce = uuid.New().String()
	_, err = s.service.Login(login, nonce, ip)
	a.Nil(err)
}

func Test_NoLogin(t *testing.T) {
	a, s := Setup(t)

	// Add noise
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()

	login := api.Login{Provider: "moovtest", SubjectID: uuid.New().String()}
	nonce := uuid.New().String()
	ip := "1.2.3.4"

	// first login should work
	_, err := s.service.Login(login, nonce, ip)
	a.NotNil(err)
}

func Test_DisableNonExistantCredential(t *testing.T) {
	a, s := Setup(t)

	// Add noise
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()
	_, _ = s.RegisterRandom()

	_, err := s.service.DisableCredentials(s.session, uuid.New().String(), uuid.New().String())
	a.NotNil(err)
}
