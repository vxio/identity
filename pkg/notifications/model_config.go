package notifications

type NotificationsConfig struct {
	SMTP *SMTPConfig
}

type SMTPConfig struct {
	host string
	port int
	user string
	pass string
	from string
}
