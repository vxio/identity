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

func Test_CustomKeyValue(t *testing.T) {
	a, buffer, log := Setup(t)

	log.WithKeyValue("custom", "value").Log("test")

	a.Contains(buffer.String(), "custom=value")
}

func Test_CustomMap(t *testing.T) {
	a, buffer, log := Setup(t)

	log.WithMap(map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}).Log("test")

	output := buffer.String()
	a.Contains(output, "custom1=value1")
	a.Contains(output, "custom2=value2")
}

func Test_MultipleContexts(t *testing.T) {
	a, buffer, log := Setup(t)

	log.
		WithKeyValue("custom1", "value1").
		WithKeyValue("custom2", "value2").
		Log("test")

	output := buffer.String()
	a.Contains(output, "custom1=value1")
	a.Contains(output, "custom2=value2")
}

func Setup(t *testing.T) (*assert.Assertions, *strings.Builder, Logger) {
	a := assert.New(t)

	buffer := strings.Builder{}
	writer := log.NewLogfmtLogger(log.NewSyncWriter(&buffer))
	log := NewLogger(writer)

	return a, &buffer, log
}
