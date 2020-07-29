package invites

import (
	"testing"
	"time"

	"github.com/gorilla/mux"
	api "github.com/moov-io/identity/pkg/api"
	authntestutils "github.com/moov-io/identity/pkg/authn/testutils"
	client "github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	tmwt "github.com/moov-io/tumbler/pkg/middleware/middlewaretest"
)

type Scope struct {
	session       tmw.TumblerClaims
	config        Config
	time          stime.StaticTimeService
	notifications notifications.NotificationsService
	repository    Repository
	service       api.InvitesApiServicer
	routes        *mux.Router
	api           *client.APIClient
}

func NewScope(t *testing.T) Scope {
	logging := logging.NewDefaultLogger()
	session := tmwt.NewRandomClaims()

	invitesConfig := Config{
		Expiration: time.Hour,
		SendToHost: "http://local.moov.io",
		SendToPath: "http://local.moov.io",
	}

	times := stime.NewStaticTimeService()
	repository := NewInMemoryInvitesRepository(t)

	notifications := notifications.NewMockNotificationsService(notifications.MockConfig{
		From: "noreply@moov.io",
	})

	authnClient := authntestutils.NewMockAuthnClient()

	service, err := NewInvitesService(invitesConfig, times, repository, notifications, authnClient)
	if err != nil {
		t.Error(err)
	}

	controller := NewInvitesController(service)

	routes := mux.NewRouter()
	api.AppendRouters(logging, routes, controller)

	testMiddleware := tmwt.NewTestMiddleware(times, session)
	routes.Use(testMiddleware.Handler)

	testAPI := clienttest.NewTestClient(routes)

	return Scope{
		session:       session,
		config:        invitesConfig,
		time:          times,
		notifications: notifications,
		repository:    repository,
		service:       service,
		routes:        routes,
		api:           testAPI,
	}
}
