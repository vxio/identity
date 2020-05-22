package authn

import (
	"testing"
	"time"

	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	fuzz "github.com/google/gofuzz"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"github.com/moov-io/identity/pkg/zerotrust"
	"github.com/stretchr/testify/require"
)

func Setup(t *testing.T) (*require.Assertions, Scope, *fuzz.Fuzzer) {
	a := require.New(t)

	logger := log.NewNopLogger()
	session := zerotrust.NewRandomSession()

	db, close, err := database.NewAndMigrate(database.InMemorySqliteConfig, logger, nil)
	t.Cleanup(close)
	a.Nil(err)

	stime := stime.NewStaticTimeService()

	notifications := notifications.NewMockNotificationsService(notifications.MockConfig{From: "noreply@moov.io"})

	invitesConfig := invites.Config{
		Expiration: time.Hour,
		SendToURL:  "https://localhost/register",
	}
	invitesRepo := invites.NewInvitesRepository(db)
	invites, err := invites.NewInvitesService(invitesConfig, stime, invitesRepo, notifications)
	a.Nil(err)

	identitiesRepo := identities.NewIdentityRepository(db)
	identities := identities.NewIdentitiesService(stime, identitiesRepo)

	credsRepo := credentials.NewCredentialRepository(db)
	creds := credentials.NewCredentialsService(stime, credsRepo)

	identityKeys, err := webkeys.NewGenerateJwksService()
	a.Nil(err)

	sessionConfig := SessionConfig{
		Expiration: time.Hour,
		LandingURL: "https://localhost/whoami",
	}
	token := NewSessionService(stime, identityKeys, sessionConfig)

	service := NewAuthnService(*creds, *identities, token, invites, sessionConfig.LandingURL)

	f := fuzz.New().Funcs(
		func(e *LoginSession, c fuzz.Continue) {
			e.IP = "1.2.3.4"
			e.State = c.RandString()

			e.StandardClaims = jwt.StandardClaims{
				ExpiresAt: stime.Now().Add(time.Hour).Unix(),
				NotBefore: stime.Now().Unix(),
				IssuedAt:  stime.Now().Unix(),
				Id:        uuid.New().String(),
				Subject:   uuid.New().String(),
				Audience:  "moovauth",
				Issuer:    "moovauth",
			}

			e.Register = client.Register{
				Provider:   c.RandString(),
				SubjectID:  uuid.New().String(),
				InviteCode: c.RandString(),
				Email:      c.RandString() + "@moovtest.io",
			}
		},
	)

	return a, Scope{
		session: session,
		stime:   stime,
		logger:  logger,
		service: service,
		invites: invites,
	}, f
}

type Scope struct {
	session zerotrust.Session
	stime   stime.StaticTimeService
	logger  log.Logger
	service api.InternalApiServicer
	invites api.InvitesApiServicer
}

func (s *Scope) NewClient(loginSession LoginSession) *client.APIClient {
	testAuthnMiddleware := NewTestMiddleware(s.stime, loginSession)
	controller := NewAuthnAPIController(s.logger, s.service)

	routes := mux.NewRouter()
	api.AppendRouters(routes, controller)
	routes.Use(testAuthnMiddleware.Handler)

	testAPI := clienttest.NewTestClient(routes)

	return testAPI
}

// TestMiddleware - Handles injecting a session into a request for testing
type TestMiddleware struct {
	time    stime.TimeService
	session LoginSession
}

// NewTestMiddleware - Generates a default Middleware that always injects the specified Session into the request
func NewTestMiddleware(time stime.TimeService, session LoginSession) *TestMiddleware {
	return &TestMiddleware{
		time:    time,
		session: session,
	}
}

// Handler - Generates the handler you use to wrap the http routes
func (s *TestMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't really like using this map of any objects in the context for this, but it seems how its done.
		ctx := context.WithValue(r.Context(), LoginSessionContextKey, &s.session)

		h.ServeHTTP(w, r.Clone(ctx))
	})
}
