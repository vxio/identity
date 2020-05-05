package jwks

import "time"

type JwksConfig struct {
	File *JwksFileConfig
	HTTP *JwksHttpConfig

	Expiration time.Duration
}

type JwksFileConfig struct {
	Path string
}

type JwksHttpConfig struct {
	URL string
}
