package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier_configuration"

	"github.com/google/uuid"
)

type Repository interface {
	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	CreateIdentifierConfiguration(ctx context.Context, configuration *identifier_configuration.IdentifierConfiguration) error
	UpdateIdentifierConfiguration(ctx context.Context, id uuid.UUID, updateFn func(config *identifier_configuration.IdentifierConfiguration) (*identifier_configuration.IdentifierConfiguration, error)) error

	CreateIdentifier(ctx context.Context, identifier *identifier.Identifier) error

	Migrate(ctx context.Context) error
}
