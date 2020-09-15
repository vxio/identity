package session_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/session"
)

func Test_Generate_Cookie(t *testing.T) {
	s := NewSessionScope(t)

	r, err := http.NewRequest("GET", "http://local.moov.io", nil)
	s.assert.Nil(err)

	cookie, err := s.token.GenerateCookie(r, session.Session{
		CredentialID: uuid.New(),
		IdentityID:   uuid.New(),
		TenantID:     uuid.New(),
	})
	s.assert.Nil(err)

	s.assert.Equal("moov", cookie.Name)
}
