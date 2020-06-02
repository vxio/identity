package invites

import (
	"time"
)

// Config holds the configuration for the Invites package
type Config struct {
	Expiration time.Duration
	SendToURL  string
}
