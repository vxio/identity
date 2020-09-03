package notifications

import (
	"crypto/tls"
	"fmt"

	log "github.com/moov-io/identity/pkg/logging"
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
	d.TLSConfig = &tls.Config{
		ServerName:         config.Host,
		InsecureSkipVerify: config.InsecureSSL,
	}

	return &smtpService{
		logger:    logger,
		dailer:    d,
		config:    config,
		templates: templates,
	}
}

func (s *smtpService) SendEmail(to string, email EmailTemplate) error {
	// Never log the message because of security concerns.
	logCtx := s.logger.WithMap(map[string]string{
		"email_to":       to,
		"email_from":     s.config.From,
		"email_subject":  email.EmailSubject(),
		"email_template": email.TemplateName(),
	})

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", email.EmailSubject())

	txt, err := s.templates.Text(email)
	if err != nil {
		return logCtx.Error().LogError("Unable to generate text template", err)
	}

	if txt != "" {
		m.SetBody("text/plain", txt)
	}

	html, err := s.templates.HTML(email)
	if err != nil {
		return logCtx.Error().LogError("Unable to generate html template", err)
	}

	if html != "" {
		m.SetBody("text/html", html)
	}

	if err := s.dailer.DialAndSend(m); err != nil {
		return logCtx.Error().LogError("Failed to send email", err)
	}

	logCtx.Info().Log(fmt.Sprintf("Successfully sent email to: %s", to))
	return nil
}
