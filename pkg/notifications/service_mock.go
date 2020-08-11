package notifications

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
	s.sent = append(s.sent, email)
	return nil
}
