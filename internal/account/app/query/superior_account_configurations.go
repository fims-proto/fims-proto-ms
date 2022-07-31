package query

import (
	"context"
	"github.com/google/uuid"
)

type SuperiorAccountConfigurationsReadModel interface {
	SuperiorAccountConfigurations(ctx context.Context, accountId uuid.UUID) ([]AccountConfiguration, error)
}
