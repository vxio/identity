package notifications

import authn "github.com/moov-io/authn/pkg/client"

type InviteEmail struct {
	Subject             string
	AcceptInvitationURL string
	Tenant              authn.Tenant
}

func NewInviteEmail(url string, tenant authn.Tenant) InviteEmail {
	return InviteEmail{
		Subject:             "Invite for Moov.io!",
		AcceptInvitationURL: url,
		Tenant:              tenant,
	}
}

func (i *InviteEmail) TemplateName() string {
	return "invite.template"
}

func (i *InviteEmail) EmailSubject() string {
	return i.Subject
}
