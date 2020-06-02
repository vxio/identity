package session

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Session struct {
	IdentityID   uuid.UUID `json:"iid"`
	TenantID     uuid.UUID `json:"tid"`
	CredentialID uuid.UUID `json:"cid"`
}

type SessionJwt struct {
	jwt.StandardClaims

	Session
}
