package jwks

type JwksConfig struct {
	File *JwksFileConfig
	HTTP *JwksHttpConfig
}

type JwksFileConfig struct {
	Path string
}

type JwksHttpConfig struct {
	URL string
}
