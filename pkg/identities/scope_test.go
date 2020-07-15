package identities_test

import (
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/database"
	. "github.com/moov-io/identity/pkg/identities"
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
		service:    *service,
		api:        testAPI,
	}
}

func Setup(t *testing.T) (*require.Assertions, Scope, *fuzz.Fuzzer) {
	a := require.New(t)
	s := NewScope(t)
	f := NewFuzzer()
	return a, s, f
}

func (s *Scope) RandomInvite() api.Invite {
	return api.Invite{
		InviteID: uuid.New().String(),
		TenantID: s.session.TenantID.String(),
	}
}

func RandStringN(c fuzz.Continue, n int) string {
	s := make([]rune, n)
	for i := range s {
		c.Fuzz(&s[i])
	}
	return string(s)
}

func NewFuzzer() *fuzz.Fuzzer {
	return fuzz.New().NumElements(0, 5).Funcs(
		func(e *api.RegisterAddress, c fuzz.Continue) {
			e.Type = "primary"
			e.Address1 = c.RandString()
			c.Fuzz(e.Address2)
			e.City = c.RandString()
			e.State = RandStringN(c, 2)
			e.PostalCode = RandStringN(c, 5)
			e.Country = RandStringN(c, 2)
		},

		func(e *client.UpdateAddress, c fuzz.Continue) {
			e.Type = "primary"
			e.Address1 = c.RandString()
			c.Fuzz(e.Address2)
			e.City = c.RandString()
			e.State = RandStringN(c, 2)
			e.PostalCode = RandStringN(c, 5)
			e.Country = RandStringN(c, 2)
			e.Validated = c.RandBool()
		},

		func(e *api.Address, c fuzz.Continue) {
			e.Type = "primary"
			e.Address1 = c.RandString()
			c.Fuzz(e.Address2)
			e.City = c.RandString()
			e.State = RandStringN(c, 2)
			e.PostalCode = RandStringN(c, 5)
			e.Country = RandStringN(c, 2)
		},

		func(e *api.RegisterPhone, c fuzz.Continue) {
			e.Type = "mobile"
			e.Number = RandStringN(c, 15)
		},

		func(e *client.UpdatePhone, c fuzz.Continue) {
			e.Type = "mobile"
			e.Number = RandStringN(c, 15)
			e.Validated = c.RandBool()
		},

		func(e *api.Phone, c fuzz.Continue) {
			e.Type = "mobile"
			e.Number = RandStringN(c, 15)
		},

		func(e *api.Register, c fuzz.Continue) {
			e.CredentialID = uuid.New().String()

			e.InviteCode = c.RandString()
			e.FirstName = c.RandString()
			e.MiddleName = c.RandString()
			e.LastName = c.RandString()
			c.Fuzz(e.NickName)
			c.Fuzz(e.Suffix)
			e.BirthDate = time.Unix(c.Int63n(time.Now().Unix()), 0).In(time.UTC)
			e.Email = c.RandString() + "@test.moov.io"

			e.Phones = make([]api.RegisterPhone, c.Intn(3)+1)
			for i := range e.Phones {
				c.Fuzz(&e.Phones[i])
			}

			e.Addresses = make([]api.RegisterAddress, c.Intn(3)+1)
			for i := range e.Addresses {
				c.Fuzz(&e.Addresses[i])
			}
		},

		func(e *api.UpdateIdentity, c fuzz.Continue) {
			e.FirstName = c.RandString()
			e.MiddleName = c.RandString()
			e.LastName = c.RandString()
			c.Fuzz(e.NickName)
			c.Fuzz(e.Suffix)
			e.BirthDate = time.Unix(c.Int63n(time.Now().Unix()), 0).In(time.UTC)
			e.Status = RandStringN(c, 10)

			e.Phones = make([]client.UpdatePhone, c.Intn(3)+1)
			for i := range e.Phones {
				c.Fuzz(&e.Phones[i])
			}

			e.Addresses = make([]client.UpdateAddress, c.Intn(3)+1)
			for i := range e.Addresses {
				c.Fuzz(&e.Addresses[i])
			}
		},
	)
}
