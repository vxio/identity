package gateway

import (
	"github.com/google/uuid"
)

type TenantID uuid.UUID

func (id *TenantID) String() string {
	return uuid.UUID(*id).String()
}
