package notifications

type NotificationsConfig struct {
	SMTP *SMTPConfig
	Mock *MockConfig
}

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
	SSL  bool
}

type MockConfig struct {
	From string
}
