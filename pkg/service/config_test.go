package service_test

import (
	"testing"

	configpkg "github.com/moov-io/identity/pkg/config"
	log "github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/service"
	"github.com/stretchr/testify/require"
)

func Test_ConfigLoading(t *testing.T) {
	logger := log.NewNopLogger()

	ConfigService := configpkg.NewConfigService(logger)

	config := &service.GlobalConfig{}
	err := ConfigService.Load(config)
	require.Nil(t, err)
}
