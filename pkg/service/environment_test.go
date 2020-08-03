package service_test

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/service"
	"github.com/stretchr/testify/assert"
)

func Test_Environment_Startup(t *testing.T) {
	a := assert.New(t)

	logger := logging.NewLogger(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)))
	env, err := service.NewEnvironment(logger, nil)
	a.Nil(err)

	shutdown := env.RunServers(false)
	t.Cleanup(shutdown)
}
