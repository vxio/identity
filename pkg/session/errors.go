package session

import "errors"

var (
	ErrIdentityNotFound     = errors.New("identity not set or found")
	ErrCredentialsNotSet    = errors.New("credentialID not set")
	ErrPutSessionNotEnabled = errors.New("put session not enabled")
)
