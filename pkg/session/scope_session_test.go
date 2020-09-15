package session_test

import (
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/tumbler/pkg/jwe"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	tmwt "github.com/moov-io/tumbler/pkg/middleware/middlewaretest"
	"github.com/moov-io/tumbler/pkg/webkeys"
	"github.com/stretchr/testify/require"
)

type SessionScope struct {
	t          *testing.T
	assert     *require.Assertions
	claims     tmw.TumblerClaims
	time       stime.StaticTimeService
	identities identities.Service
	controller session.SessionController
	token      session.TokenService
}

func NewSessionScope(t *testing.T) SessionScope {
	a := require.New(t)

	claims := tmwt.NewRandomClaims()
	logging := logging.NewDefaultLogger()
	times := stime.NewStaticTimeService()

	config := session.Config{
		Expiration:       time.Hour,
		EnablePutSession: true,
	}

	db, close, err := database.NewAndMigrate(database.InMemorySqliteConfig, nil, nil)
	t.Cleanup(close)
	a.Nil(err)

	keys, err := webkeys.NewGenerateJwksService()
	a.Nil(err)

	jwe := jwe.NewJWEService(times, time.Hour, keys)

	credentialsRepo := credentials.NewCredentialRepository(db)
	credentials := credentials.NewCredentialsService(times, credentialsRepo)

	identitiesRepository := identities.NewIdentityRepository(db)
	identities := identities.NewIdentitiesService(times, identitiesRepository)

	token := session.NewTokenService(times, jwe, config)
	service := session.NewSessionService(logging, identities, token, credentials, config)

	controller := session.NewSessionController(logging, service)

	return SessionScope{
		t:          t,
		assert:     a,
		claims:     claims,
		time:       times,
		identities: identities,
		controller: controller,
		token:      token,
	}
}

func (s *SessionScope) APIClient() *client.APIClient {
	routes := mux.NewRouter()
	s.controller.AppendRoutes(routes)

	testMiddleware := tmwt.NewTestMiddleware(s.time, s.claims)
	routes.Use(testMiddleware.Handler)

	testAPI := clienttest.NewTestClient(routes)
	return testAPI
}
