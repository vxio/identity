package notifications

type NotificationsService interface {
	SendInvite(email string, secretCode string, acceptInvitationUrl string) error
}

func NewNotificationsService(config NotificationsConfig) NotificationsService {
	if config.SMTP != nil {
		return NewSmtpNotificationsService(*config.SMTP)
	}

	return nil
}
