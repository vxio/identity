package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/log"

	"github.com/markbates/pkger"
	"github.com/spf13/viper"
)

type ConfigService struct {
	logger log.Logger
	path   string
}

func NewConfigService(logger log.Logger) ConfigService {
	return ConfigService{
		logger: logger,
	}
}

func (s *ConfigService) Load(config interface{}) error {
	s.logger.Log("config", "Loading config")

	err := s.LoadFile(pkger.Include("/configs/config.default.yml"), config)
	if err != nil {
		return err
	}

	if file, ok := os.LookupEnv("APP_CONFIG"); ok && strings.TrimSpace(file) != "" {
		s.logger.Log("config", fmt.Sprintf("Loading config - %s", file))
		overrides := viper.New()
		overrides.SetConfigFile(file)

		if err := overrides.ReadInConfig(); err != nil {
			msg := fmt.Sprintf("Failed loading the specific app config - %s", err)
			s.logger.Log("config", msg)
			return err
		}

		if err := overrides.Unmarshal(config); err != nil {
			msg := fmt.Sprintf("Unable to unmarshal the specific app config - %s", err)
			s.logger.Log("config", msg)
			return err
		}
	}

	return nil
}

func (s *ConfigService) LoadFile(file string, config interface{}) error {
	s.logger.Log("config", "Loading config", "file", file)

	f, err := pkger.Open(file)
	if err != nil {
		s.logger.Log("config", fmt.Sprintf("Pkger unable to load %s - cause: %s", file, err.Error()))
		return err
	}

	deflt := viper.New()
	deflt.SetConfigType("yaml")
	if err := deflt.ReadConfig(f); err != nil {
		msg := "Unable to load the defaults"
		s.logger.Log("config", msg)
		return errors.New(msg)
	}

	if err := deflt.Unmarshal(config); err != nil {
		msg := fmt.Sprintf("Unable to unmarshal the defaults - %v", err)
		s.logger.Log("config", msg)
		return errors.New(msg)
	}

	return nil
}
