package config

import (
	"errors"
	"fmt"
	"os"

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

	f, err := pkger.Open("/configs/config.default.yml")
	if err != nil {
		s.logger.Log("config", fmt.Sprintf("Pkger unable to load config.default.yml - cause: %s", err.Error()))
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

	overrides := viper.New()

	if v, s := os.LookupEnv("APP_CONFIG"); s {
		overrides.SetConfigFile(v)
	} else {
		overrides.SetConfigFile("./configs/config.overrides.yml")
	}

	if err := overrides.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); ok {
			s.logger.Log("config", fmt.Sprintf("Didn't find override config, just continuing - %v", err))
			return nil
		} else {
			msg := fmt.Sprintf("Failed loading the override - %v", err)
			s.logger.Log("config", msg)
			return err
		}
	}

	if err := overrides.Unmarshal(config); err != nil {
		return errors.New(fmt.Sprintf("Unable to load overrides config - %s", err))
	}

	return nil
}
