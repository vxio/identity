package notifications

import (
	"crypto/tls"
	"fmt"

	"github.com/go-kit/kit/log"
	"gopkg.in/gomail.v2"
)

type smtpService struct {
	logger    log.Logger
	dailer    gomail.Dialer
	config    SMTPConfig
	templates TemplateRepository
}

func NewSmtpNotificationsService(logger log.Logger, config SMTPConfig, templates TemplateRepository) NotificationsService {
	d := *gomail.NewDialer(config.Host, config.Port, config.User, config.Pass)
	d.SSL = config.SSL
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &smtpService{
		logger:    logger,
		dailer:    d,
		config:    config,
		templates: templates,
	}
}

func (s *smtpService) SendEmail(to string, email EmailTemplate) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", email.EmailSubject())

	txt, err := s.templates.Text(email)
	if err != nil {
		return err
	}

	if txt != "" {
		m.SetBody("text/plain", txt)
	}

	html, err := s.templates.HTML(email)
	if err != nil {
		return err
	}

	if html != "" {
		m.SetBody("text/html", html)
	}

	if err := s.dailer.DialAndSend(m); err != nil {
		s.logger.Log("level", "error", "msg", fmt.Sprintf("Failed to send email - %s\n", err))
		return err
	}

	s.logger.Log("level", "info", "msg", fmt.Sprintf("Successfully sent email to: %s", to))
	return nil
}
