package notifications

import (
	authn "github.com/moov-io/authn/pkg/client"
	"github.com/moov-io/identity/pkg/client"
)

type InviteEmail struct {
	Subject             string
	AcceptInvitationURL string
	Inviter             client.Identity
	Tenant              authn.Tenant
}

func NewInviteEmail(url string, inviter client.Identity, tenant authn.Tenant) InviteEmail {
	return InviteEmail{
		Subject:             "Invite for Moov.io!",
		AcceptInvitationURL: url,
		Inviter:             inviter,
		Tenant:              tenant,
	}
}

func (i *InviteEmail) TemplateName() string {
	return "invite.template"
}

func (i *InviteEmail) EmailSubject() string {
	return i.Subject
}
