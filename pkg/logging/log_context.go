package logging

type LogContext interface {
	LogContext() map[string]string
}
