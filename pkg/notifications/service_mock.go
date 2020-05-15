package notifications

import (
	"fmt"
)

type mockService struct {
	config MockConfig
}

func NewMockNotificationsService(config MockConfig) NotificationsService {
	return &mockService{
		config: config,
	}
}

func (s *mockService) SendInvite(email string, secretCode string, acceptInvitationUrl string) error {
	subject := "Invite for Moov.io"
	body := "Here is your invite for Moov.io"

	fmt.Printf("  From: %s\n  To: %s\n  Subject: %s\n  Message: %s\n", s.config.From, email, subject, body)

	return nil
}
