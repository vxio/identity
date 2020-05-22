package authn

import (
	"context"
	"testing"

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
	_, resp, err := c.InternalApi.RegisterWithCredentials(context.Background(), "asdf", ls.Register)
	a.Nil(err)
	a.Equal(302, resp.StatusCode)
	
	asfdfas come back here and fix this asdfasdfvar

	//client.InternalApi.RegisterWithCredentials(context.Background())
}
