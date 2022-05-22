package domain

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type User struct {
	id     uuid.UUID
	traits json.RawMessage
}

func NewUser(id uuid.UUID, traits json.RawMessage) (*User, error) {
	if id == uuid.Nil {
		return nil, errors.New("empty user id")
	}
	if len(traits) == 0 {
		return nil, errors.New("empty user traits")
	}

	return &User{
		id:     id,
		traits: traits,
	}, nil
}

func (u User) Id() uuid.UUID {
	return u.id
}

func (u User) Traits() json.RawMessage {
	return u.traits
}
