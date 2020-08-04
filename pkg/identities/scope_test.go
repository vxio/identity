package identities_test

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/database"
	. "github.com/moov-io/identity/pkg/identities"
	identitiestestutils "github.com/moov-io/identity/pkg/identities/testutils"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	tmwt "github.com/moov-io/tumbler/pkg/middleware/middlewaretest"
	"github.com/stretchr/testify/require"
)

type Scope struct {
	session    tmw.TumblerClaims
	time       stime.StaticTimeService
	repository Repository
	service    Service
	api        *client.APIClient
}

func NewScope(t *testing.T) Scope {
	logging := logging.NewDefaultLogger()
	session := tmwt.NewRandomClaims()
	times := stime.NewStaticTimeService()

	db, close, err := database.NewAndMigrate(database.InMemorySqliteConfig, nil, nil)
	t.Cleanup(close)
	if err != nil {
		t.Error(err)
	}

	repository := NewIdentityRepository(db)

	service := NewIdentitiesService(times, repository)

	controller := NewIdentitiesController(service)

	routes := mux.NewRouter()
	api.AppendRouters(logging, routes, controller)

	testMiddleware := tmwt.NewTestMiddleware(times, session)
	routes.Use(testMiddleware.Handler)

	testAPI := clienttest.NewTestClient(routes)

	return Scope{
		session:    session,
		time:       times,
		repository: repository,
		service:    service,
		api:        testAPI,
	}
}

func Setup(t *testing.T) (*require.Assertions, Scope, *fuzz.Fuzzer) {
	a := require.New(t)
	s := NewScope(t)
	f := identitiestestutils.NewFuzzer()
	return a, s, f
}

func (s *Scope) RandomInvite() client.Invite {
	return client.Invite{
		InviteID: uuid.New().String(),
		TenantID: s.session.TenantID.String(),
	}
}
