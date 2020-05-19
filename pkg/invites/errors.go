package invites

import "errors"

var ErrTokenExpired = errors.New("Invite token is expired")
var ErrTokenDisabled = errors.New("Invite was disabled")
