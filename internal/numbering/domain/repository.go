package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateIdentifierConfiguration(ctx context.Context, configuration *IdentifierConfiguration) error
	UpdateIdentifierConfiguration(ctx context.Context, id uuid.UUID, updateFn func(config *IdentifierConfiguration) (*IdentifierConfiguration, error)) error

	CreateIdentifier(ctx context.Context, identifier *Identifier) error

	Migrate(ctx context.Context) error
}
