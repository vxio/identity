package jwks

import (
	"github.com/go-kit/kit/log"

	"gopkg.in/square/go-jose.v2"
)

type JwksService interface {
	FetchJwks() (*jose.JSONWebKeySet, error)
}

func NewJwksService(logger log.Logger, config JwksConfig) (JwksService, error) {
	if config.File != nil {
		return NewFileJwksService(config.File.Path), nil
	} else if config.HTTP != nil {
		return NewWebJwksService(config.HTTP.URL), nil
	} else {
		logger.Log("jwks", "Generating new JWKS keys")
		return NewGenerateJwksService()
	}
}
