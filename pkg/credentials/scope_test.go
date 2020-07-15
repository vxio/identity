package credentials_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	. "github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	tmwt "github.com/moov-io/tumbler/pkg/middleware/middlewaretest"
	"github.com/stretchr/testify/require"
)

type Scope struct {
	session    tmw.TumblerClaims
	time       stime.StaticTimeService
	repository CredentialRepository
	service    CredentialsService
	api        *client.APIClient
}

func NewScope(t *testing.T) Scope {
	logger := logging.NewDefaultLogger()
	session := tmwt.NewRandomClaims()
	times := stime.NewStaticTimeService()

	db, close, err := database.NewAndMigrate(database.InMemorySqliteConfig, nil, nil)
	t.Cleanup(close)
	if err != nil {
		t.Error(err)
	}

	repository := NewCredentialRepository(db)

	service := NewCredentialsService(times, repository)

	controller := NewCredentialsApiController(service)

	routes := mux.NewRouter()
	api.AppendRouters(logger, routes, controller)

	testMiddleware := tmwt.NewTestMiddleware(times, session)
	routes.Use(testMiddleware.Handler)

	testAPI := clienttest.NewTestClient(routes)

	return Scope{
		session:    session,
		time:       times,
		repository: repository,
		service:    *service,
		api:        testAPI,
	}
}

func Setup(t *testing.T) (*require.Assertions, Scope) {
	a := require.New(t)
	s := NewScope(t)
	return a, s
}

func (s *Scope) RegisterRandom() (*api.Credential, error) {
	identityID := uuid.New().String()
	credentialID := uuid.New().String()

	return s.service.Register(identityID, credentialID, s.session.TenantID.String())
}
