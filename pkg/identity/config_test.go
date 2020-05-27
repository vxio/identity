package identity_test

import (
	"testing"

	"github.com/go-kit/kit/log"
	configpkg "github.com/moov-io/identity/pkg/config"
	. "github.com/moov-io/identity/pkg/identity"
	"github.com/stretchr/testify/require"
)

func Test_ConfigLoading(t *testing.T) {
	logger := log.NewNopLogger()

	ConfigService := configpkg.NewConfigService(logger)

	config := &Config{}
	err := ConfigService.Load(config)
	require.Nil(t, err)
}
