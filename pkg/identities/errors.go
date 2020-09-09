package identities

import "errors"

var (
	ErrIdentityNotFound = errors.New("no identity found for that tenant")
)
