package authn_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/moov-io/identity/pkg/authn"
	"github.com/moov-io/identity/pkg/client"
)

func Test_Register(t *testing.T) {
	s := Setup(t)

	invite, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	s.assert.Nil(err)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.TenantID = invite.TenantID
	ls.InviteCode = code
	ls.Scopes = []string{"register", "finished"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.Equal(200, resp.StatusCode)
	s.assert.Nil(err)
}

func Test_Register_InvalidInviteCode(t *testing.T) {
	s := Setup(t)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.TenantID = s.session.TenantID.String()
	ls.InviteCode = "doesnotexist"
	ls.Scopes = []string{"register", "finished"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.NotNil(err)
	s.assert.Equal(400, resp.StatusCode)
}

func Test_Register_Invalid_Scope(t *testing.T) {
	s := Setup(t)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.TenantID = s.session.TenantID.String()
	ls.InviteCode = "doesnotexist"
	ls.Scopes = []string{"badscope"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.NotNil(err)
	s.assert.Equal(404, resp.StatusCode)
}

func Test_Login_Failed(t *testing.T) {
	s := Setup(t)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.Scopes = []string{"authenticate", "finished"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.Authenticated(context.Background())
	s.assert.Equal(404, resp.StatusCode)
	s.assert.NotNil(err)
}

func Test_Login_Invalid_Scope(t *testing.T) {
	s := Setup(t)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.Scopes = []string{"badscope"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.Authenticated(context.Background())
	s.assert.Equal(404, resp.StatusCode)
	s.assert.NotNil(err)
}

func Test_Login_Success(t *testing.T) {
	s := Setup(t)

	registerSession := RegisterRandomIdentity(s)

	loginSession := LoginSession{}
	s.fuzz.Fuzz(&loginSession)

	// These are the values that have to match up to what was registered.
	loginSession.CredentialID = registerSession.CredentialID
	loginSession.TenantID = registerSession.TenantID
	loginSession.Scopes = []string{"authenticate", "finished"}

	// Test if we can login with it.
	c := s.NewClient(loginSession)
	_, resp, err := c.AuthenticationApi.Authenticated(context.Background())
	s.assert.Equal(200, resp.StatusCode)
	s.assert.Nil(err)
}

func RegisterRandomIdentity(s Scope) LoginSession {
	req := httptest.NewRequest("GET", "https://local.moov.io", strings.NewReader(""))
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)

	// First need to invite the user
	invite, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	if err != nil {
		panic(err)
	}

	// Register the user with the code.
	registerSession := LoginSession{}
	s.fuzz.Fuzz(&registerSession)
	registerSession.TenantID = invite.TenantID
	registerSession.InviteCode = code
	registerSession.Scopes = []string{"register"}

	_, err = s.service.RegisterWithCredentials(req, registerSession.Register, registerSession.State, registerSession.IP)
	if err != nil {
		panic(err)
	}

	return registerSession
}
