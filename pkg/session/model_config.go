package session

import (
	"time"

	"github.com/moov-io/identity/pkg/webkeys"
)

// Config - Holds the configuration for the session cookie created after registration or logging in
type Config struct {
	Expiration  time.Duration
	PublicKeys  webkeys.WebKeysConfig
	PrivateKeys webkeys.WebKeysConfig
}
