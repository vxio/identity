package gateway

import (
	"encoding/json"

	"github.com/google/uuid"
)

type TenantID uuid.UUID

func (id *TenantID) String() string {
	return uuid.UUID(*id).String()
}

func (id *TenantID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(*id))
}

func (id *TenantID) UnmarshalJSON(data []byte) error {
	parsed, err := uuid.ParseBytes(data)
	if err != nil {
		return err
	}

	*id = TenantID(parsed)
	return nil
}
