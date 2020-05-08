package webkeys

import (
	"github.com/go-kit/kit/log"

	"gopkg.in/square/go-jose.v2"
)

type WebKeysService interface {
	FetchJwks() (*jose.JSONWebKeySet, error)
}

func NewWebKeysService(logger log.Logger, config WebKeysConfig) (WebKeysService, error) {
	if config.File != nil {
		return NewFileJwksService(config.File.Path), nil
	} else if config.HTTP != nil {
		return NewHTTPJwksService(config.HTTP.URL), nil
	} else {
		logger.Log("jwks", "Generating new JWKS keys")
		return NewGenerateJwksService()
	}
}
