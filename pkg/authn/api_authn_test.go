package authn_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
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
	if err != nil {
		fmt.Printf("%s", s.logOutput.String())
	}
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)
}

func Test_Register_TwoTenants(t *testing.T) {
	s := Setup(t)

	credential := uuid.New().String()

	ls := LoginSession{}

	// Do the first tenant
	s.fuzz.Fuzz(&ls)
	ls.CredentialID = credential
	ls.TenantID = uuid.New().String()
	ls.Scopes = []string{"register", "finished", "signup"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	if err != nil {
		fmt.Printf("%s", s.logOutput.String())
	}
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)

	// Lets do the second one.
	s.fuzz.Fuzz(&ls)
	ls.CredentialID = credential
	ls.TenantID = uuid.New().String()
	ls.Scopes = []string{"register", "finished", "signup"}

	c = s.NewClient(ls)
	_, resp, err = c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)
}

func Test_Register_TwiceSameTenant(t *testing.T) {
	s := Setup(t)

	credential := uuid.New().String()
	tenant := uuid.New().String()

	ls := LoginSession{}

	// Do the first tenant
	s.fuzz.Fuzz(&ls)
	ls.CredentialID = credential
	ls.TenantID = tenant
	ls.Scopes = []string{"register", "finished", "signup"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	if err != nil {
		fmt.Printf("%s", s.logOutput.String())
	}
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)

	// Lets do the second one.
	s.fuzz.Fuzz(&ls)
	ls.CredentialID = credential
	ls.TenantID = tenant
	ls.Scopes = []string{"register", "finished", "signup"}

	c = s.NewClient(ls)
	_, resp, err = c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.NotNil(err)
	s.assert.Equal(404, resp.StatusCode)
}

func Test_Register_Success_Returns_ImageURL_If_Available(t *testing.T) {
	s := Setup(t)

	invite, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	s.assert.Nil(err)

	ls := LoginSession{}

	s.fuzz.Fuzz(&ls)
	imageUrl := "http://images.com/123.jpg"
	ls.Register.ImageUrl = &imageUrl
	ls.TenantID = invite.TenantID
	ls.InviteCode = code
	ls.Scopes = []string{"register", "finished"}

	c := s.NewClient(ls)
	loggedIn, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.Equal(ls.Register.ImageUrl, loggedIn.ImageUrl)
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)
}

func Test_Register_Success_Returns_Empty_ImageURL(t *testing.T) {
	s := Setup(t)

	invite, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	s.assert.Nil(err)

	ls := LoginSession{}

	s.fuzz.Fuzz(&ls)
	ls.TenantID = invite.TenantID
	ls.InviteCode = code
	ls.Scopes = []string{"register", "finished"}

	c := s.NewClient(ls)
	loggedIn, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.Equal(ls.Register.ImageUrl, loggedIn.ImageUrl)
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)
}

func Test_Signup(t *testing.T) {
	s := Setup(t)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.TenantID = s.session.TenantID.String()
	ls.Scopes = []string{"register", "finished", "signup"}

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), ls.Register)
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)
}

func Test_Register_PartialJson(t *testing.T) {
	s := Setup(t)

	invite, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	s.assert.Nil(err)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.TenantID = invite.TenantID
	ls.InviteCode = code
	ls.Scopes = []string{"register", "finished"}

	c := s.NewClient(ls)
	loggedIn, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), client.Register{})
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)

	s.assert.Equal(ls.FirstName, loggedIn.FirstName)
	s.assert.Equal(ls.LastName, loggedIn.LastName)
}

func Test_Register_PartialJsonOverride(t *testing.T) {
	s := Setup(t)

	invite, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	s.assert.Nil(err)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.TenantID = invite.TenantID
	ls.InviteCode = code
	ls.Scopes = []string{"register", "finished"}

	c := s.NewClient(ls)
	loggedIn, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), client.Register{
		LastName: "overrode",
	})
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)

	s.assert.Equal(ls.FirstName, loggedIn.FirstName)
	s.assert.Equal("overrode", loggedIn.LastName)
}

func Test_Register_FailIfSessionMissingData(t *testing.T) {
	s := Setup(t)

	invite, code, err := s.invites.SendInvite(s.session, client.SendInvite{Email: "test@moovtest.io"})
	s.assert.Nil(err)

	ls := LoginSession{}
	s.fuzz.Fuzz(&ls)
	ls.TenantID = invite.TenantID
	ls.InviteCode = code
	ls.Scopes = []string{"register", "finished"}

	// Make the first name invalid so it should fail until we fix it by passing in a new one.
	ls.FirstName = ""

	c := s.NewClient(ls)

	// Lets call without passing in a new one to notify if we need them to correct any data we collected.
	loggedIn, resp, err := c.AuthenticationApi.RegisterWithCredentials(context.Background(), client.Register{})
	s.assert.NotNil(err)
	s.assert.Equal(400, resp.StatusCode)

	i, ok := err.(client.GenericOpenAPIError)
	s.assert.True(ok)

	regErrs, ok := i.Model().(client.RegisterErrors)
	s.assert.True(ok)

	// Check that the error is for the first name
	s.assert.NotEmpty(regErrs.FirstName)

	// Lets call again with fixed data
	loggedIn, resp, err = c.AuthenticationApi.RegisterWithCredentials(context.Background(), client.Register{
		FirstName: "John Doe",
	})
	s.assert.Nil(err)
	s.assert.Equal(200, resp.StatusCode)

	s.assert.Equal(200, resp.StatusCode)
	s.assert.Equal("John Doe", loggedIn.FirstName)
	s.assert.Equal(ls.LastName, loggedIn.LastName)
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
	s.assert.Equal(404, resp.StatusCode)
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
	ls.TenantID = s.session.TenantID.String()

	c := s.NewClient(ls)
	_, resp, err := c.AuthenticationApi.Authenticated(context.Background())
	s.assert.Equal(404, resp.StatusCode)
	s.assert.NotNil(err)
}
func Test_Login_MissingTenant(t *testing.T) {
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

	_, _, err = s.service.RegisterWithCredentials(req, registerSession.Register, registerSession.State, registerSession.IP, false)
	if err != nil {
		panic(err)
	}

	return registerSession
}
