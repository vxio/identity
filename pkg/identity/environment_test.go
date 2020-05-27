package identity_test

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/moov-io/identity/pkg/identity"
	"github.com/stretchr/testify/assert"
)

func Test_Environment_Startup(t *testing.T) {
	a := assert.New(t)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	env, err := identity.NewEnvironment(logger, nil)
	a.Nil(err)

	shutdown := env.RunServers(false)
	t.Cleanup(shutdown)
}
