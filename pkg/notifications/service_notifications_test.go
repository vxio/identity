package notifications

import (
	"testing"

	"github.com/google/uuid"
	"github.com/moov-io/base/docker"
	"github.com/moov-io/identity/pkg/authn"
	authntestutils "github.com/moov-io/identity/pkg/authn/testutils"
	log "github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/tumbler/pkg/middleware"
	"github.com/moov-io/tumbler/pkg/middleware/middlewaretest"
	"github.com/stretchr/testify/assert"
)

func Test_SMTP_SendInvite(t *testing.T) {
	if !docker.Enabled() {
		t.SkipNow()
	}

	a, s := Setup(t)

	config := NotificationsConfig{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        2025,
			User:        "test",
			Pass:        "test",
			From:        "noreply@moovtest.io",
			SSL:         true,
			InsecureSSL: true,
		},
	}

	service, err := NewNotificationsService(s.logger, config, s.templates)
	a.Nil(err, "Check that `docker-compose up` is running before running tests. Can't talk to mailslurper.")

	tenant, err := s.authn.GetTenant(s.claims, uuid.New().String())
	a.Nil(err)

	invite := NewInviteEmail("https://localhost/accept", *tenant)

	err = service.SendEmail("test@moovtest.io", &invite)
	a.Nil(err, "Check that `docker-compose up` is running before running tests. Can't talk to mailslurper.")
}

func Test_Mock_SendInvite(t *testing.T) {
	a, s := Setup(t)

	config := NotificationsConfig{
		Mock: &MockConfig{
			From: "noreply@moovtest.io",
		},
	}

	service, err := NewNotificationsService(s.logger, config, s.templates)
	a.Nil(err, "Check that `docker-compose up` is running before running tests. Can't talk to mailslurper.")

	tenant, err := s.authn.GetTenant(s.claims, uuid.New().String())
	a.Nil(err)

	invite := NewInviteEmail("https://localhost/accept", *tenant)

	err = service.SendEmail("test@moovtest.io", &invite)
	a.Nil(err, "Check that `docker-compose up` is running before running tests. Can't talk to mailslurper.")

	mock, ok := service.(*mockService)
	a.True(ok)

	a.Contains(mock.sent, &invite)

}

type Scope struct {
	logger    log.Logger
	templates TemplateRepository
	authn     authn.AuthnClient
	claims    middleware.TumblerClaims
}

func Setup(t *testing.T) (*assert.Assertions, Scope) {
	a := assert.New(t)

	authn := authntestutils.NewMockAuthnClient()

	logger := log.NewNopLogger()
	templateRepository, err := NewTemplateRepository(logger)
	a.Nil(err)

	return a, Scope{
		logger:    logger,
		templates: templateRepository,
		authn:     authn,
		claims:    middlewaretest.NewRandomClaims(),
	}
}
