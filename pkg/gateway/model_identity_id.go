package gateway

import (
	"encoding/json"

	"github.com/google/uuid"
)

type IdentityID uuid.UUID

func (id *IdentityID) String() string {
	return uuid.UUID(*id).String()
}

func (id *IdentityID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(*id))
}

func (id *IdentityID) UnmarshalJSON(data []byte) error {
	parsed, err := uuid.ParseBytes(data)
	if err != nil {
		return err
	}

	*id = IdentityID(parsed)
	return nil
}
