package authn

import (
	"time"
)

type SessionConfig struct {
	Expiration time.Duration
	LandingURL string
}
