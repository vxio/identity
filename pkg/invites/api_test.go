package invites

import (
	"context"
	"testing"

	"github.com/moov-io/identity/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestAPIInvite(t *testing.T) {
	// New assertion so we don't have to pass t in on every call
	a := assert.New(t)
	s := NewScope(t)

	invite := sendInvite(a, s, "newinvite@moov.io")

	a.Equal(s.session.Subject, invite.InvitedBy, "Invited By doesn't match session caller")
	a.Equal(s.time.Now(), invite.InvitedOn)
	a.Equal(s.session.TenantID.String(), invite.TenantID, "TenantID doesn't match session tenantID")

	a.Equal(s.time.Now().Add(s.config.Expiration), invite.ExpiresOn)

	a.Nil(invite.DisabledBy)
	a.Nil(invite.DisabledOn)
	a.Nil(invite.RedeemedOn)
}

func TestAPIList(t *testing.T) {
	// New assertion so we don't have to pass t in on every call
	a := assert.New(t)
	s := NewScope(t)

	sent := sendInvite(a, s, "list@moov.io")

	invites := listInvites(a, s)

	a.Len(invites, 1)
	a.Equal(sent.InviteID, invites[0].InviteID, "InviteID's dont match")
	a.Equal(s.session.Subject, invites[0].InvitedBy, "Invited By doesn't match session caller")
	a.Equal(s.session.TenantID.String(), invites[0].TenantID, "TenantID doesn't match session tenantID")
}

func TestAPIDeactivate(t *testing.T) {
	// New assertion so we don't have to pass t in on every call
	a := assert.New(t)
	s := NewScope(t)

	sent := sendInvite(a, s, "todeactivate@moov.io")

	resp, err := s.api.InvitesApi.DisableInvite(context.Background(), sent.InviteID)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	list := listInvites(a, s)

	a.Len(list, 1)

	deactivated := list[0]
	a.NotNil(deactivated.DisabledBy)
	a.NotNil(deactivated.DisabledOn)
	a.Equal(s.session.Subject, *deactivated.DisabledBy)
	a.Equal(s.time.Now(), *deactivated.DisabledOn)
	a.Equal(sent.InviteID, deactivated.InviteID)
}

func sendInvite(a *assert.Assertions, s Scope, email string) client.Invite {
	invite, _, err := s.api.InvitesApi.SendInvite(context.Background(), client.SendInvite{Email: email})
	a.Nil(err)
	a.Equal(email, invite.Email)
	return invite
}

func listInvites(a *assert.Assertions, s Scope) []client.Invite {
	invites, _, err := s.api.InvitesApi.ListInvites(context.Background(), nil)
	a.Nil(err)
	return invites
}
