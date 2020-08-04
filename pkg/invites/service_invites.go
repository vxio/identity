package invites

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	authnclient "github.com/moov-io/identity/pkg/authn/client"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/stime"

	api "github.com/moov-io/identity/pkg/api"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// InvitesApiServicer defines the api actions for the InvitesApi service
type InvitesService interface {
	DisableInvite(tmw.TumblerClaims, string) error
	ListInvites(tmw.TumblerClaims) ([]client.Invite, error)
	SendInvite(tmw.TumblerClaims, client.SendInvite) (*client.Invite, string, error)
	Redeem(code string) (*client.Invite, error)
}

type invitesService struct {
	sendToURL     *template.Template
	expiration    time.Duration
	time          stime.TimeService
	repository    Repository
	notifications notifications.NotificationsService
	authnClient   authnclient.AuthnClient
	identity      identities.Service
}

// NewInvitesService instantiates a new invitesService for interacting with Invites from outside of the package.
func NewInvitesService(config Config, time stime.TimeService, repository Repository, notifications notifications.NotificationsService, authnClient authnclient.AuthnClient, identity identities.Service) (InvitesService, error) {

	urlTemplate, err := template.New("send").Parse(config.SendToHost + config.SendToPath)
	if err != nil {
		return nil, err
	}

	return &invitesService{
		sendToURL:     urlTemplate,
		expiration:    config.Expiration,
		time:          time,
		repository:    repository,
		notifications: notifications,
		authnClient:   authnClient,
		identity:      identity,
	}, nil
}

// ListInvites - List outstanding invites
func (s *invitesService) ListInvites(claims tmw.TumblerClaims) ([]client.Invite, error) {
	invites, err := s.repository.list(api.TenantID(claims.TenantID))
	return invites, err
}

// SendInvite - Send an email invite to a new user
func (s *invitesService) SendInvite(claims tmw.TumblerClaims, send client.SendInvite) (*client.Invite, string, error) {
	invite := client.Invite{
		InviteID:   uuid.New().String(),
		TenantID:   claims.TenantID.String(),
		Email:      send.Email,
		InvitedBy:  claims.IdentityID.String(),
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

	redeemURL, err := generateRedeemURL(*s.sendToURL, invite, code)
	if err != nil {
		return nil, "", err
	}

	tenant, err := s.authnClient.GetTenant(claims, claims.TenantID.String())
	if err != nil {
		return nil, "", err
	}

	inviter, err := s.identity.GetIdentity(claims, claims.IdentityID.String())
	if err != nil {
		return nil, "", err
	}

	notification := notifications.NewInviteEmail(redeemURL.String(), *inviter, *tenant)

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
func (s *invitesService) DisableInvite(claims tmw.TumblerClaims, inviteID string) error {
	invite, err := s.repository.get(api.TenantID(claims.TenantID), inviteID)
	if err != nil {
		return err
	}

	disabledBy := claims.IdentityID.String()
	disabledOn := s.time.Now()
	invite.DisabledBy = &disabledBy
	invite.DisabledOn = &disabledOn

	return s.repository.update(*invite)
}

func (s *invitesService) Redeem(code string) (*client.Invite, error) {
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

func generateRedeemURL(sendToURL template.Template, invite client.Invite, code *string) (*url.URL, error) {
	data := struct {
		TenantID string
	}{invite.TenantID}

	urlString := strings.Builder{}
	err := sendToURL.Execute(&urlString, data)
	if err != nil {
		return nil, err
	}

	// duplicate it so we can append the invite code to the mutable value
	sendTo, err := url.Parse(urlString.String())
	if err != nil {
		return nil, err
	}
	qry := sendTo.Query()
	qry.Add("invite_code", *code)
	sendTo.RawQuery = qry.Encode()

	return sendTo, nil
}
