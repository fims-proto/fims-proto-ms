package identifier

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Identifier struct {
	id                        uuid.UUID
	identifierConfigurationId uuid.UUID
	identifier                string
}

func New(id, identifierConfigurationId uuid.UUID, identifier string) (*Identifier, error) {
	if id == uuid.Nil {
		return nil, errors.New("id cannot be empty")
	}

	if identifierConfigurationId == uuid.Nil {
		return nil, errors.New("identifier configuration id cannot be empty")
	}

	if identifier == "" {
		return nil, errors.New("identifier cannot be empty")
	}

	return &Identifier{
		id:                        id,
		identifierConfigurationId: identifierConfigurationId,
		identifier:                identifier,
	}, nil
}

func (i Identifier) Id() uuid.UUID {
	return i.id
}

func (i Identifier) IdentifierConfigurationId() uuid.UUID {
	return i.identifierConfigurationId
}

func (i Identifier) Identifier() string {
	return i.identifier
}
