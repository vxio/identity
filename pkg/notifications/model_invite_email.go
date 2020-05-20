package notifications

type InviteEmail struct {
	Subject             string
	AcceptInvitationURL string
}

func NewInviteEmail(url string) InviteEmail {
	return InviteEmail{
		Subject:             "Invite for Moov.io!",
		AcceptInvitationURL: url,
	}
}

func (i *InviteEmail) TemplateName() string {
	return "invite.template"
}

func (i *InviteEmail) EmailSubject() string {
	return i.Subject
}
