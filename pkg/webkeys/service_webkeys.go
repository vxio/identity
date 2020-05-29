package webkeys

import (
	"github.com/go-kit/kit/log"

	"gopkg.in/square/go-jose.v2"
)

type WebKeysService interface {
	Keys() (*jose.JSONWebKeySet, error)
}

func NewWebKeysService(logger log.Logger, config WebKeysConfig) (WebKeysService, error) {
	if config.HTTP != nil {
		return NewHTTPJwksService(*config.HTTP, nil)
	} else if config.File != nil {
		return NewFileJwksService(*config.File)
	} else {
		logger.Log("jwks", "Generating new JWKS keys")
		return NewGenerateJwksService()
	}
}
