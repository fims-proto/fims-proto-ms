package identifier

import (
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

type Identifier struct {
	id                        uuid.UUID
	identifierConfigurationId uuid.UUID
	identifier                string
}

func New(id, identifierConfigurationId uuid.UUID, identifier string) (*Identifier, error) {
	if id == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugNumberingIdEmpty)
	}

	if identifierConfigurationId == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugNumberingConfigIdEmpty)
	}

	if identifier == "" {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugNumberingIdentifierEmpty)
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
