package invites

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/stime"
)

type InvitesServiceScope struct {
	session gateway.Session
	service api.InvitesApiServicer
	time    stime.StaticTimeService
}

func TestSendInvite(t *testing.T) {
	s := NewInvitesScope(t)

	sendInvite := api.SendInvite{Email: "testuser@moov.io"}

	invite, code, err := s.service.SendInvite(s.session, sendInvite)
	if err != nil {
		t.Error(err)
	}

	invites, err := s.service.ListInvites(s.session)
	if err != nil {
		t.Error(err)
	}

	if len(invites) != 1 {
		t.Errorf("Length of invites isn't 1")
	}

	if *invite != invites[0] {
		t.Errorf("Invite doesn't exist in list %s", cmp.Diff(*invite, invites[0]))
	}

	redeemed, err := s.service.Redeem(code)
	if err != nil {
		t.Error(err)
	}

	if redeemed.InviteID != invite.InviteID {
		t.Errorf("Redeemed InviteID doesn't match Sent InviteID")
	}
}
func Test_RedeemExpired(t *testing.T) {
	s := NewInvitesScope(t)

	sendInvite := api.SendInvite{Email: "testuser@moov.io"}

	_, _, err := s.service.SendInvite(s.session, sendInvite)
	if err != nil {
		t.Error(err)
	}

	_, err = s.service.Redeem("doesnotexist")
	if err != sql.ErrNoRows {
		t.Error("A token that does not exist didn't fail with No Rows")
	}
}

func TestDisableInvite(t *testing.T) {
	s := NewInvitesScope(t)

	sendInvite := api.SendInvite{Email: "testuser@moov.io"}

	invite, code, err := s.service.SendInvite(s.session, sendInvite)
	if err != nil {
		t.Error(err)
	}

	err = s.service.DisableInvite(s.session, invite.InviteID)
	if err != nil {
		t.Error(err)
	}

	_, err = s.service.Redeem(code)
	if err != ErrInviteCodeDisabled {
		t.Error("Disabled token didn't redeem will disabled failure")
	}
}
func TestExpiredInvite(t *testing.T) {
	s := NewInvitesScope(t)

	sendInvite := api.SendInvite{Email: "testuser@moov.io"}

	invite, code, err := s.service.SendInvite(s.session, sendInvite)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(s.time.Now())
	s.time.Change(invite.ExpiresOn.Add(time.Millisecond))
	fmt.Println(s.time.Now())

	_, err = s.service.Redeem(code)
	if err != ErrInviteCodeExpired {
		t.Error("Expired token didn't redeem with expired failure")
	}
}

func NewInvitesScope(t *testing.T) InvitesServiceScope {
	session := gateway.NewRandomSession()

	repository := NewInMemoryInvitesRepository(t)

	config := Config{
		Expiration: time.Hour,
		SendToURL:  "http://local.moov.io",
	}

	times := stime.NewStaticTimeService()
	notification := notifications.NewMockNotificationsService(notifications.MockConfig{
		From: "noreply@moov.io",
	})

	service, err := NewInvitesService(config, times, repository, notification)
	if err != nil {
		panic(err)
	}

	scope := InvitesServiceScope{
		session: session,
		service: service,
		time:    times,
	}

	return scope
}
