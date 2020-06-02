package authn

import "github.com/moov-io/identity/pkg/webkeys"

type Config struct {
	LandingURL string
	Keys       webkeys.WebKeysConfig
}
