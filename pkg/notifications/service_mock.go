package notifications

import (
	"fmt"
)

type mockService struct {
	config MockConfig
	sent   []EmailTemplate
}

func NewMockNotificationsService(config MockConfig) NotificationsService {
	return &mockService{
		config: config,
	}
}

func (s *mockService) SendEmail(to string, email EmailTemplate) error {
	subject := email.EmailSubject()

	s.sent = append(s.sent, email)

	fmt.Printf("  From: %s\n  To: %s\n  Subject: %s\n  Template: %+v\n", s.config.From, email, subject, email)

	return nil
}
