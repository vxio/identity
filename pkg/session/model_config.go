package session

import (
	"time"
)

// Config - Holds the configuration for the session cookie created after registration or logging in
type Config struct {
	Expiration time.Duration
}
