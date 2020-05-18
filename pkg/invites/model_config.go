package invites

import (
	"time"
)

type InvitesConfig struct {
	Expiration time.Duration
	SendToURL  string
}
