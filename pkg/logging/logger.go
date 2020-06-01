package logging

import (
	"fmt"

	"github.com/go-kit/kit/log"
)

type LogContext interface {
	LogContext() map[string]string
}

type Logger interface {
	With(ctxs ...LogContext) Logger
	Log(msg string)

	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

type logger struct {
	writer log.Logger
	ctx    map[string]string
}

func NewLogger(writer log.Logger) Logger {
	return &logger{
		writer: writer,
		ctx:    map[string]string{},
	}
}

// With returns a new Logger with the contexts added to its own.
func (l *logger) With(ctxs ...LogContext) Logger {
	// Estimation assuming that for each ctxs has at least 1 value.
	combined := make(map[string]string, len(l.ctx)+len(ctxs))

	for k, v := range l.ctx {
		combined[k] = v
	}

	for _, c := range ctxs {
		itemCtx := c.LogContext()
		for k, v := range itemCtx {
			combined[k] = v
		}
	}

	return &logger{
		writer: l.writer,
		ctx:    combined,
	}
}

func (l *logger) Log(msg string) {
	i := 0
	keyvals := make([]interface{}, (len(l.ctx)*2)+2)
	for k, v := range l.ctx {
		keyvals[i] = k
		keyvals[i+1] = v
		i += 2
	}

	keyvals[i] = "msg"
	keyvals[i+1] = msg

	fmt.Println(keyvals)

	l.writer.Log(keyvals...)
}

func (l *logger) Info(msg string) {
	l.With(Info).Log(msg)
}

func (l *logger) Error(msg string) {
	l.With(Error).Log(msg)
}

func (l *logger) Fatal(msg string) {
	l.With(Fatal).Log(msg)
}
