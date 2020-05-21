package invites

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	api "github.com/moov-io/identity/pkg/api"
	client "github.com/moov-io/identity/pkg/client"
	clienttest "github.com/moov-io/identity/pkg/client_test"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/zerotrust"
	"github.com/moov-io/identity/pkg/zerotrust/zerotrusttest"
)

type Scope struct {
	session       zerotrust.Session
	config        Config
	time          stime.StaticTimeService
	notifications notifications.NotificationsService
	repository    Repository
	service       api.InvitesApiServicer
	routes        *mux.Router
	api           *client.APIClient
}

func NewScope(t *testing.T) Scope {
	session := NewSession()

	invitesConfig := Config{
		Expiration: time.Hour,
		SendToURL:  "http://local.moov.io",
	}

	times := stime.NewStaticTimeService()
	repository := NewInMemoryInvitesRepository(t)

	notifications := notifications.NewMockNotificationsService(notifications.MockConfig{
		From: "noreply@moov.io",
	})

	service, err := NewInvitesService(invitesConfig, times, repository, notifications)
	if err != nil {
		t.Error(err)
	}

	controller := NewInvitesController(service)

	routes := mux.NewRouter()
	api.AppendRouters(routes, controller)

	testMiddleware := zerotrusttest.NewTestMiddleware(times, session)
	routes.Use(testMiddleware.Handler)

	testApi := clienttest.NewTestClient(routes)

	return Scope{
		session:       session,
		config:        invitesConfig,
		time:          times,
		notifications: notifications,
		repository:    repository,
		service:       service,
		routes:        routes,
		api:           testApi,
	}
}

func NewSession() zerotrust.Session {
	return zerotrust.Session{
		CallerID: zerotrust.IdentityID(uuid.New()),
		TenantID: zerotrust.TenantID(uuid.New()),
	}
}
