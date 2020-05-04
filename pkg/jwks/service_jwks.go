package jwks

import (
	"gopkg.in/square/go-jose.v2"
)

type JwksService interface {
	FetchJwks() (*jose.JSONWebKeySet, error)
}
