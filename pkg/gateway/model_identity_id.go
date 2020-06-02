package gateway

import (
	"github.com/google/uuid"
)

type IdentityID uuid.UUID

func (id *IdentityID) String() string {
	return uuid.UUID(*id).String()
}
