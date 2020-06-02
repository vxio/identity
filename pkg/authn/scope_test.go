package authn_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	. "github.com/moov-io/identity/pkg/authn"
	log "github.com/moov-io/identity/pkg/logging"

	"github.com/dgrijalva/jwt-go"
	fuzz "github.com/google/gofuzz"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/notifications"
	sessionpkg "github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"github.com/stretchr/testify/require"
)

// pull these out so it speeds up testing
var authnKeys, _ = webkeys.NewGenerateJwksService()
var identityKeys, _ = webkeys.NewGenerateJwksService()

func Setup(t *testing.T) (*require.Assertions, Scope, *fuzz.Fuzzer) {
	a := require.New(t)

	logger := log.NewDefaultLogger()
	session := gateway.NewRandomSession()

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

	sessionConfig := sessionpkg.Config{Expiration: time.Hour}
	token := sessionpkg.NewSessionService(stime, identityKeys, sessionConfig)

	authnConfig := Config{LandingURL: "https://localhost/whoami"}
	service := NewAuthnService(logger, *creds, *identities, token, invites, authnConfig.LandingURL)

	f := fuzz.New().Funcs(
		func(e *LoginSession, c fuzz.Continue) {
			e.IP = "1.2.3.4"
			e.State = c.RandString()

			e.StandardClaims = jwt.StandardClaims{
				ExpiresAt: stime.Now().Add(time.Hour).Unix(),
				NotBefore: stime.Now().Add(time.Minute * -1).Unix(),
				IssuedAt:  stime.Now().Add(time.Minute * -1).Unix(),
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
		sessionConfig: sessionConfig,
		authnConfig:   authnConfig,
		session:       session,
		stime:         stime,
		logger:        logger,
		service:       service,
		invites:       invites,
		authnKeys:     authnKeys,
	}, f
}

type Scope struct {
	sessionConfig sessionpkg.Config
	authnConfig   Config
	session       gateway.Session
	stime         stime.StaticTimeService
	logger        log.Logger
	service       api.InternalApiServicer
	invites       api.InvitesApiServicer
	authnKeys     *webkeys.GenerateJwksService
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

func (s *Scope) Cookie(session LoginSession) *http.Cookie {
	privateKey := s.authnKeys.Private
	signingMethod := jwt.GetSigningMethod(privateKey.Algorithm)

	token := jwt.NewWithClaims(signingMethod, session)
	token.Header["kid"] = privateKey.KeyID

	tokenString, err := token.SignedString(privateKey.Key)
	if err != nil {
		panic(err)
	}

	return &http.Cookie{
		Name:     "moov-authn",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Unix(session.ExpiresAt, 0),
		MaxAge:   int(time.Unix(session.ExpiresAt, 0).Second()),
		SameSite: http.SameSiteDefaultMode,
		Secure:   false,
		HttpOnly: true,
	}
}
