package invites

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/stime"

	api "github.com/moov-io/identity/pkg/api"
)

type invitesService struct {
	sendToURL     url.URL
	expiration    time.Duration
	time          stime.TimeService
	repository    Repository
	notifications notifications.NotificationsService
}

// NewInvitesService instantiates a new invitesService for interacting with Invites from outside of the package.
func NewInvitesService(config Config, time stime.TimeService, repository Repository, notifications notifications.NotificationsService) (api.InvitesApiServicer, error) {

	url, err := url.Parse(config.SendToURL)
	if err != nil {
		return nil, err
	}

	return &invitesService{
		sendToURL:     *url,
		expiration:    config.Expiration,
		time:          time,
		repository:    repository,
		notifications: notifications,
	}, nil
}

// ListInvites - List outstanding invites
func (s *invitesService) ListInvites(session gateway.Session) ([]api.Invite, error) {
	invites, err := s.repository.list(session.TenantID)
	return invites, err
}

// SendInvite - Send an email invite to a new user
func (s *invitesService) SendInvite(session gateway.Session, send api.SendInvite) (*api.Invite, string, error) {
	invite := api.Invite{
		InviteID:   uuid.New().String(),
		TenantID:   session.TenantID.String(),
		Email:      send.Email,
		InvitedBy:  session.CallerID.String(),
		InvitedOn:  s.time.Now(),
		RedeemedOn: nil,
		ExpiresOn:  s.time.Now().Add(s.expiration),
		DisabledBy: nil,
		DisabledOn: nil,
	}

	code, err1 := generateInviteCode()
	if err1 != nil {
		return nil, "", err1
	}

	// duplicate it so we can append the invite code to the mutable value
	sendTo, _ := url.Parse(s.sendToURL.String())
	qry := sendTo.Query()
	qry.Add("invite_code", *code)
	sendTo.RawQuery = qry.Encode()

	notification := notifications.NewInviteEmail(sendTo.String())

	if err := s.notifications.SendEmail(invite.Email, &notification); err != nil {
		return nil, "", err
	}

	// add to DB
	created, err2 := s.repository.add(invite, *code)
	if err2 != nil {
		return nil, "", err2
	}

	return created, *code, nil
}

// DeleteInvite - Delete an invite that was sent and invalidate the token.
func (s *invitesService) DisableInvite(session gateway.Session, inviteID string) error {
	invite, err := s.repository.get(session.TenantID, inviteID)
	if err != nil {
		return err
	}

	disabledBy := session.CallerID.String()
	disabledOn := s.time.Now()
	invite.DisabledBy = &disabledBy
	invite.DisabledOn = &disabledOn

	return s.repository.update(*invite)
}

func (s *invitesService) Redeem(code string) (*api.Invite, error) {
	invite, err := s.repository.getByCode(strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}

	if invite.ExpiresOn.Before(s.time.Now()) {
		return nil, ErrInviteCodeExpired
	}

	if invite.DisabledOn != nil {
		return nil, ErrInviteCodeDisabled
	}

	redeemedOn := s.time.Now()
	invite.RedeemedOn = &redeemedOn

	if err := s.repository.update(*invite); err != nil {
		return nil, err
	}

	return invite, nil
}

// Generate a large random crypto string to work as the invitation token
func generateInviteCode() (*string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	code := base64.RawURLEncoding.EncodeToString(b)
	return &code, nil
}
