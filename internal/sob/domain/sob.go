package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Sob struct {
	uuid uuid.UUID
	name string
}

func NewSob(sobUUID uuid.UUID, name string) (*Sob, error) {
	if sobUUID == uuid.Nil {
		return nil, errors.New("empty sob uuid")
	}
	if name == "" {
		return nil, errors.New("empty sob name")
	}

	return &Sob{
		uuid: sobUUID,
		name: name,
	}, nil
}

func (s Sob) UUID() uuid.UUID {
	return s.uuid
}

func (s Sob) Name() string {
	return s.name
}
