package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/moov-io/identity/pkg/logging"

	"github.com/markbates/pkger"
	"github.com/spf13/viper"
)

type ConfigService struct {
	logger logging.Logger
	path   string
}

func NewConfigService(logger logging.Logger) ConfigService {
	return ConfigService{
		logger: logger.WithKeyValue("component", "ConfigService"),
	}
}

func (s *ConfigService) Load(config interface{}) error {
	err := s.LoadFile(pkger.Include("/configs/config.default.yml"), config)
	if err != nil {
		return err
	}

	if file, ok := os.LookupEnv("APP_CONFIG"); ok && strings.TrimSpace(file) != "" {
		log := s.logger.WithKeyValue("app_config", file)
		log.Info().Log("Loading APP_CONFIG config file")

		overrides := viper.New()
		overrides.SetConfigFile(file)

		if err := overrides.ReadInConfig(); err != nil {
			return log.LogError(fmt.Sprintf("Failed loading the specific app config - %s", err), err)
		}

		if err := overrides.Unmarshal(config); err != nil {
			return log.LogError(fmt.Sprintf("Unable to unmarshal the specific app config - %s", err), err)
		}
	}

	return nil
}

func (s *ConfigService) LoadFile(file string, config interface{}) error {
	log := s.logger.WithKeyValue("file", file)
	log.Info().Log("Loading config file")

	f, err := pkger.Open(file)
	if err != nil {
		return log.LogError("Pkger unable to load", err)
	}

	deflt := viper.New()
	deflt.SetConfigType("yaml")
	if err := deflt.ReadConfig(f); err != nil {
		return log.LogError("Unable to load the defaults", err)
	}

	if err := deflt.Unmarshal(config); err != nil {
		return log.LogError(fmt.Sprintf("Unable to unmarshal the defaults"), err)
	}

	return nil
}
