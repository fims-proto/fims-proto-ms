package query

import (
	"context"
	"github.com/google/uuid"
)

type AllAccountConfigurationsReadModel interface {
	AllAccountConfigurations(ctx context.Context, sobId uuid.UUID) ([]AccountConfiguration, error)
}
