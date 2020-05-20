package authn

import (
	"time"
)

// SessionConfig - Holds the configuration for the session cookie created after registration or logging in
type SessionConfig struct {
	Expiration time.Duration
	LandingURL string
}
