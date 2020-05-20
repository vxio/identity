package notifications

import (
	"errors"

	"github.com/go-kit/kit/log"
)

type NotificationsService interface {
	SendEmail(to string, email EmailTemplate) error
}

func NewNotificationsService(logger log.Logger, config NotificationsConfig, templates TemplateRepository) (NotificationsService, error) {
	if config.SMTP != nil {
		return NewSmtpNotificationsService(logger, *config.SMTP, templates), nil
	} else if config.Mock != nil {
		return NewMockNotificationsService(*config.Mock), nil
	}

	return nil, errors.New("No notifications method specified. Check config.")
}
