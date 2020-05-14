package notifications

type NotificationsConfig struct {
	SMTP *SMTPConfig
	Mock *MockConfig
}

type SMTPConfig struct {
	host string
	port int
	user string
	pass string
	from string
}

type MockConfig struct {
	from string
}
