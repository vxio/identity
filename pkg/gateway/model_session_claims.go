package gateway

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type SessionClaims struct {
	jwt.StandardClaims

	CallerID uuid.UUID `json:"iid"`
	TenantID uuid.UUID `json:"tid"`
}
