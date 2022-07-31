package query

import (
	"context"
	"github.com/google/uuid"
)

type AccountsInPeriodReadModel interface {
	AccountsInPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) ([]Account, error)
}
