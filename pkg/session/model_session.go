package session

import (
	"github.com/google/uuid"
	"github.com/moov-io/tumbler/pkg/jwe"
)

type Session struct {
	IdentityID   uuid.UUID `json:"iid"`
	TenantID     uuid.UUID `json:"tid"`
	CredentialID uuid.UUID `json:"cid"`
}

type SessionJwt struct {
	jwe.Claims

	Session
}
