/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package identityserver

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/utils"
)

var (
	INVITE_TIMEOUT = time.Hour * 48
)

type InvitesService struct {
	time          utils.TimeService
	repository    InvitesRepository
	notifications notifications.NotificationsService
}

func NewInvitesService(time utils.TimeService, repository InvitesRepository, notifications notifications.NotificationsService) InvitesApiServicer {
	return &InvitesService{
		time:          time,
		repository:    repository,
		notifications: notifications,
	}
}

// DeleteInvite - Delete an invite that was sent and invalidate the token.
func (s *InvitesService) DeleteInvite(inviteID string) (interface{}, error) {
	s.repository.delete("1", inviteID)
	return nil, errors.New("service method 'DeleteInvite' not implemented")
}

// ListInvites - List outstanding invites
func (s *InvitesService) ListInvites(orgID string) (interface{}, error) {
	invites, err := s.repository.list("tenantID")
	return invites, err
}

// SendInvite - Send an email invite to a new user
func (s *InvitesService) SendInvite(tenant TenantID, send SendInvite) (interface{}, error) {
	invite := Invite{
		InviteID:  uuid.New().String(),
		TenantID:  string(tenant),
		Email:     send.Email,
		InvitedBy: "1", // @TODO
		InvitedOn: s.time.Now(),
		ExpiresOn: s.time.Now().Add(INVITE_TIMEOUT),
		Redeemed:  false,
	}

	code, err1 := generateInviteCode()
	if err1 != nil {
		return nil, err1
	}

	// add to DB
	created, err2 := s.repository.add(invite, *code) // @TODO tenantID
	if err2 != nil {
		return nil, err2
	}

	// send email
	if err := s.notifications.SendInvite(created.Email, *code, "someurl"); err != nil {
		// clear out the one we added to the DB so it isn't just sitting around being unused.
		s.repository.delete(created.TenantID, created.InviteID)
		return nil, err
	}

	return created, nil
}

// Generate a large random crypto string to work as the invitation token
func generateInviteCode() (*string, error) {
	b := make([]byte, 36)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	code := base64.RawStdEncoding.EncodeToString(b)
	return &code, nil
}
