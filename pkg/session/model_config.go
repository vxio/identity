package session

import (
	"time"

	"github.com/moov-io/tumbler/pkg/webkeys"
)

// Config - Holds the configuration for the session cookie created after registration or logging in
type Config struct {
	Expiration       time.Duration
	Keys             webkeys.WebKeysConfig
	EnablePutSession bool
}
