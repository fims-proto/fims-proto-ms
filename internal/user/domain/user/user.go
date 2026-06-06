package user

import (
	"encoding/json"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

type User struct {
	id     uuid.UUID
	traits json.RawMessage
}

func New(id uuid.UUID, traits json.RawMessage) (*User, error) {
	if id == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugUserEmptyId)
	}

	return &User{
		id:     id,
		traits: traits,
	}, nil
}

func (u *User) Id() uuid.UUID {
	return u.id
}

func (u *User) Traits() json.RawMessage {
	return u.traits
}
