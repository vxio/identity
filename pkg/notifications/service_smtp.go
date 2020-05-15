package notifications

import (
	"gopkg.in/gomail.v2"
)

type smtpService struct {
	dailer gomail.Dialer
	config SMTPConfig
}

func NewSmtpNotificationsService(config SMTPConfig) NotificationsService {
	return &smtpService{
		dailer: *gomail.NewDialer(config.Host, config.Port, config.User, config.Pass),
		config: config,
	}
}

func (s *smtpService) SendInvite(email string, secretCode string, acceptInvitationUrl string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Invite for Moov.io")
	m.SetBody("text/plain", "Here is your invite for Moov.io")

	if err := s.dailer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
