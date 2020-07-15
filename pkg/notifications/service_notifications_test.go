package notifications

import (
	"testing"

	"github.com/moov-io/base/docker"
	log "github.com/moov-io/identity/pkg/logging"
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

	invite := NewInviteEmail("https://localhost/accept")

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

	invite := NewInviteEmail("https://localhost/accept")

	err = service.SendEmail("test@moovtest.io", &invite)
	a.Nil(err, "Check that `docker-compose up` is running before running tests. Can't talk to mailslurper.")

	mock, ok := service.(*mockService)
	a.True(ok)

	a.Contains(mock.sent, &invite)

}

type Scope struct {
	logger    log.Logger
	templates TemplateRepository
}

func Setup(t *testing.T) (*assert.Assertions, Scope) {
	a := assert.New(t)

	logger := log.NewNopLogger()
	templateRepository, err := NewTemplateRepository(logger)
	a.Nil(err)

	return a, Scope{
		logger:    logger,
		templates: templateRepository,
	}
}
