package webkeys

import (
	log "github.com/moov-io/identity/pkg/logging"

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
		logger.Info().Log("Generating new JWKS keys")
		return NewGenerateJwksService()
	}
}
