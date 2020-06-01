package logging

import (
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_Log(t *testing.T) {
	a, buffer, log := Setup(t)

	log.Log("my message")

	a.Contains(buffer.String(), "my message")
}

func Test_WithContext(t *testing.T) {
	a, buffer, log := Setup(t)

	log.With(Error).Log("my error message")

	a.Contains(buffer.String(), "level=error")
}

func Test_ReplaceContextValue(t *testing.T) {
	a, buffer, log := Setup(t)

	log.With(Error).With(Info).Log("my error message")

	a.Contains(buffer.String(), "level=info")
}

func Test_Info(t *testing.T) {
	a, buffer, log := Setup(t)

	log.Info("message")

	a.Contains(buffer.String(), "level=info")
}

func Test_Error(t *testing.T) {
	a, buffer, log := Setup(t)

	log.Error("message")

	a.Contains(buffer.String(), "level=error")
}

func Test_Fatal(t *testing.T) {
	a, buffer, log := Setup(t)

	log.Fatal("message")

	a.Contains(buffer.String(), "level=fatal")
}

func Setup(t *testing.T) (*assert.Assertions, *strings.Builder, Logger) {
	a := assert.New(t)

	buffer := strings.Builder{}
	writer := log.NewLogfmtLogger(log.NewSyncWriter(&buffer))
	log := NewLogger(writer)

	return a, &buffer, log
}
