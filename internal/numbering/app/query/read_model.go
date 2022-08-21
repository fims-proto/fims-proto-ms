package query

import (
	"context"

	"github.com/google/uuid"
)

type NumberingReadModel interface {
	ResolveIdentifierConfiguration(ctx context.Context, targetBusinessObject string, objectsToMatch map[string]string) (IdentifierConfiguration, error)

	IdentifierById(ctx context.Context, id uuid.UUID) (Identifier, error)
}
