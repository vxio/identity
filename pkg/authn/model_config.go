package authn

import "github.com/moov-io/tumbler/pkg/webkeys"

type Config struct {
	LandingURL string
	Keys       webkeys.WebKeysConfig
}
