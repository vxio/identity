package notifications

import "errors"

type NotificationsService interface {
	SendInvite(email string, secretCode string, acceptInvitationUrl string) error
}

func NewNotificationsService(config NotificationsConfig) (NotificationsService, error) {
	if config.SMTP != nil {
		return NewSmtpNotificationsService(*config.SMTP), nil
	} else if config.Mock != nil {
		return NewMockNotificationsService(*config.Mock), nil
	}

	return nil, errors.New("No notifications method specified. Check config.")
}
