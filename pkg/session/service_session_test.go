package session_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/tumbler/pkg/jwe"
	"github.com/moov-io/tumbler/pkg/webkeys"
	"github.com/stretchr/testify/require"
)

func Test_Generate_Cookie(t *testing.T) {
	assert := require.New(t)
	times := stime.NewStaticTimeService()
	keys, err := webkeys.NewGenerateJwksService()
	assert.Nil(err)
	jwe := jwe.NewJWEService(times, time.Hour, keys)

	service := session.NewSessionService(times, jwe, session.Config{
		Expiration: time.Hour,
	})

	r, err := http.NewRequest("GET", "http://local.moov.io", nil)
	assert.Nil(err)

	cookie, err := service.GenerateCookie(r, session.Session{
		CredentialID: uuid.New(),
		IdentityID:   uuid.New(),
		TenantID:     uuid.New(),
	})
	assert.Nil(err)

	assert.Equal("moov", cookie.Name)
}
