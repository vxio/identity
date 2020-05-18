package notifications

type Template interface {
	TemplateName() string
}

type EmailTemplate interface {
	EmailSubject() string
	Template
}
