package identity_test

import (
	"testing"

	configpkg "github.com/moov-io/identity/pkg/config"
	. "github.com/moov-io/identity/pkg/identity"
	log "github.com/moov-io/identity/pkg/logging"
	"github.com/stretchr/testify/require"
)

func Test_ConfigLoading(t *testing.T) {
	logger := log.NewNopLogger()

	ConfigService := configpkg.NewConfigService(logger)

	config := &GlobalConfig{}
	err := ConfigService.Load(config)
	require.Nil(t, err)
}
