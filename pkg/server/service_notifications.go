package identityserver

import (
	"gopkg.in/gomail.v2"
)

type NotificationsService interface {
	SendInvite(email string, secretCode string, acceptInvitationUrl string) error
}

type notificationsService struct {
	dailer gomail.Dialer
	from   string
}

func NewNotificationsService(host string, port int, user string, pass string, from string) NotificationsService {
	return &notificationsService{
		dailer: *gomail.NewDialer(host, port, user, pass),
		from:   from,
	}
}

func (s *notificationsService) SendInvite(email string, secretCode string, acceptInvitationUrl string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Invite for Moov.io")
	m.SetBody("text/plain", "Here is your invite for Moov.io")

	if err := s.dailer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
