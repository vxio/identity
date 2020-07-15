package authn_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/moov-io/identity/pkg/authn"
	log "github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/tumbler/pkg/jwe"
	"github.com/square/go-jose/jwt"

	fuzz "github.com/google/gofuzz"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moov-io/authn/pkg/keygen"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/notifications"
	sessionpkg "github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	tmwt "github.com/moov-io/tumbler/pkg/middleware/middlewaretest"
	"github.com/moov-io/tumbler/pkg/webkeys"
	"github.com/stretchr/testify/require"
)

// pull these out so it speeds up testing
var authnKeys, _ = keygen.GenerateKeys()
var identityKeys, _ = webkeys.NewGenerateJwksService()

func Setup(t *testing.T) Scope {
	a := require.New(t)

	logger := log.NewDefaultLogger()
	session := tmwt.NewRandomClaims()

	db, close, err := database.NewAndMigrate(database.InMemorySqliteConfig, logger, nil)
	t.Cleanup(close)
	a.Nil(err)

	stime := stime.NewStaticTimeService()

	notifications := notifications.NewMockNotificationsService(notifications.MockConfig{From: "noreply@moov.io"})

	invitesConfig := invites.Config{
		Expiration: time.Hour,
		SendToHost: "https://localhost",
		SendToPath: "/register",
	}
	invitesRepo := invites.NewInvitesRepository(db)
	invites, err := invites.NewInvitesService(invitesConfig, stime, invitesRepo, notifications)
	a.Nil(err)

	identitiesRepo := identities.NewIdentityRepository(db)
	identities := identities.NewIdentitiesService(stime, identitiesRepo)

	credsRepo := credentials.NewCredentialRepository(db)
	creds := credentials.NewCredentialsService(stime, credsRepo)

	sessionConfig := sessionpkg.Config{Expiration: time.Hour}
	sessionJwe := jwe.NewJWEService(stime, sessionConfig.Expiration, identityKeys)
	token := sessionpkg.NewSessionService(stime, sessionJwe, sessionConfig)

	service := authn.NewAuthnService(logger, *creds, *identities, token, invites)

	authnJwe := jwe.NewJWEService(stime, sessionConfig.Expiration, webkeys.NewStaticJwksService(authnKeys))

	f := fuzz.New().Funcs(
		func(e *authn.LoginSession, c fuzz.Continue) {
			e.IP = "1.2.3.4"
			e.State = c.RandString()

			e.Claims = jwe.Claims{
				Expiry:    jwt.NewNumericDate(stime.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(stime.Now().Add(time.Minute * -1)),
				IssuedAt:  jwt.NewNumericDate(stime.Now().Add(time.Minute * -1)),
				ID:        uuid.New().String(),
				Subject:   uuid.New().String(),
				Audience:  jwt.Audience{e.IP},
				Issuer:    "http://local.moov.io/",
			}

			e.Register = client.Register{
				CredentialID: uuid.New().String(),
				InviteCode:   c.RandString(),
				Email:        c.RandString() + "@moovtest.io",
			}
		},
	)

	return Scope{
		assert:        a,
		fuzz:          f,
		sessionConfig: sessionConfig,
		session:       session,
		stime:         stime,
		logger:        logger,
		service:       service,
		invites:       invites,
		authnJwe:      authnJwe,
		identityJwe:   sessionJwe,
	}
}

type Scope struct {
	assert        *require.Assertions
	fuzz          *fuzz.Fuzzer
	sessionConfig sessionpkg.Config
	session       tmw.TumblerClaims
	stime         stime.StaticTimeService
	logger        log.Logger
	service       api.InternalApiServicer
	invites       api.InvitesApiServicer
	authnJwe      jwe.JWEService
	identityJwe   jwe.JWEService
}

func (s *Scope) NewClient(loginSession authn.LoginSession) *client.APIClient {
	testAuthnMiddleware := NewTestMiddleware(s.stime, loginSession)

	controller := authn.NewAuthnAPIController(s.logger, s.service)

	routes := mux.NewRouter()
	api.AppendRouters(s.logger, routes, controller)
	routes.Use(testAuthnMiddleware.Handler)

	testAPI := clienttest.NewTestClient(routes)

	return testAPI
}

// TestMiddleware - Handles injecting a session into a request for testing
type TestMiddleware struct {
	time    stime.TimeService
	session authn.LoginSession
}

// NewTestMiddleware - Generates a default Middleware that always injects the specified Session into the request
func NewTestMiddleware(time stime.TimeService, session authn.LoginSession) *TestMiddleware {
	return &TestMiddleware{
		time:    time,
		session: session,
	}
}

// Handler - Generates the handler you use to wrap the http routes
func (s *TestMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't really like using this map of any objects in the context for this, but it seems how its done.
		ctx := context.WithValue(r.Context(), authn.LoginSessionContextKey, &s.session)

		h.ServeHTTP(w, r.Clone(ctx))
	})
}

func (s *Scope) AddSession(req *http.Request, modify func(session *authn.LoginSession)) authn.LoginSession {
	ls := authn.LoginSession{}
	s.fuzz.Fuzz(&ls)
	claims, err := s.authnJwe.Start(req)
	s.assert.Nil(err)
	ls.Claims = *claims

	modify(&ls)

	req.AddCookie(s.Cookie(ls))
	return ls
}

func (s *Scope) Cookie(session authn.LoginSession) *http.Cookie {
	tokenString, err := s.authnJwe.SerializeEncrypted(&session.Claims, &session)
	if err != nil {
		panic(err)
	}

	return &http.Cookie{
		Name:     "moov-authn",
		Value:    tokenString,
		Path:     "/",
		Expires:  session.Expiry.Time(),
		MaxAge:   session.Expiry.Time().Second(),
		SameSite: http.SameSiteDefaultMode,
		Secure:   false,
		HttpOnly: true,
	}
}
