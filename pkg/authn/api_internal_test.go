package authn_test

import (
	"context"
	"testing"

	fuzz "github.com/google/gofuzz"
	. "github.com/moov-io/identity/pkg/authn"
	"github.com/moov-io/identity/pkg/client"
)

func Test_Register(t *testing.T) {
	a, s, f := Setup(t)

	_, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	a.Nil(err)

	ls := LoginSession{}
	f.Fuzz(&ls)
	ls.InviteCode = code

	c := s.NewClient(ls)
	_, resp, err := c.InternalApi.RegisterWithCredentials(context.Background(), ls.Register)
	a.Equal(302, resp.StatusCode)

	redirectTo, err := resp.Location()
	a.Nil(err)
	a.Equal(redirectTo.String(), s.sessionConfig.LandingURL)
}

func Test_Register_InvalidInviteCode(t *testing.T) {
	a, s, f := Setup(t)

	ls := LoginSession{}
	f.Fuzz(&ls)
	ls.InviteCode = "doesnotexist"

	c := s.NewClient(ls)
	_, resp, err := c.InternalApi.RegisterWithCredentials(context.Background(), ls.Register)
	a.NotNil(err)
	a.Equal(400, resp.StatusCode)
}

func Test_Login_Failed(t *testing.T) {
	a, s, f := Setup(t)

	ls := LoginSession{}
	f.Fuzz(&ls)

	c := s.NewClient(ls)
	_, resp, err := c.InternalApi.Authenticated(context.Background())
	a.Equal(404, resp.StatusCode)
	a.NotNil(err)
}

func Test_Login_Success(t *testing.T) {
	a, s, f := Setup(t)

	registerSession := RegisterRandomIdentity(f, s)

	loginSession := LoginSession{}
	f.Fuzz(&loginSession)

	// These are the values that have to match up to what was registered.
	loginSession.Provider = registerSession.Provider
	loginSession.SubjectID = registerSession.SubjectID

	// Test if we can login with it.
	c := s.NewClient(loginSession)
	_, resp, err := c.InternalApi.Authenticated(context.Background())
	a.Equal(302, resp.StatusCode)
	a.NotNil(err)
}

func RegisterRandomIdentity(f *fuzz.Fuzzer, s Scope) LoginSession {

	// First need to invite the user
	_, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	if err != nil {
		panic(err)
	}

	// Register the user with the code.
	registerSession := LoginSession{}
	f.Fuzz(&registerSession)
	registerSession.InviteCode = code

	_, err = s.service.RegisterWithCredentials(registerSession.Register, registerSession.State, registerSession.IP)
	if err != nil {
		panic(err)
	}

	return registerSession
}
