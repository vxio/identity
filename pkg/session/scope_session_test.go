package session_test

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	tmwt "github.com/moov-io/tumbler/pkg/middleware/middlewaretest"
	"github.com/stretchr/testify/require"
)

type SessionScope struct {
	t          *testing.T
	assert     *require.Assertions
	claims     tmw.TumblerClaims
	time       stime.StaticTimeService
	identities identities.Service
	controller session.SessionController
}

func NewSessionScope(t *testing.T) SessionScope {
	a := require.New(t)

	claims := tmwt.NewRandomClaims()
	logging := logging.NewDefaultLogger()
	times := stime.NewStaticTimeService()

	db, close, err := database.NewAndMigrate(database.InMemorySqliteConfig, nil, nil)
	t.Cleanup(close)
	a.Nil(err)

	credRepo := credentials.NewCredentialRepository(db)
	credService := credentials.NewCredentialsService(times, credRepo)

	identitiesRepository := identities.NewIdentityRepository(db)
	identities := identities.NewIdentitiesService(times, identitiesRepository, credService)

	controller := session.NewSessionController(logging, identities, times)

	return SessionScope{
		t:          t,
		assert:     a,
		claims:     claims,
		time:       times,
		identities: identities,
		controller: controller,
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
