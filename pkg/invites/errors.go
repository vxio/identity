package invites

import "errors"

// ErrInviteCodeExpired is issued when the invite code has Expired.
var ErrInviteCodeExpired = errors.New("invite token is expired")

// ErrInviteCodeDisabled is issued when the invite was disabled by another person
var ErrInviteCodeDisabled = errors.New("invite was disabled")
