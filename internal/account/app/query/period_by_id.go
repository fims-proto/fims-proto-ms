package query

import (
	"context"
	"github.com/google/uuid"
)

type PeriodByIdReadModel interface {
	PeriodById(ctx context.Context, periodId uuid.UUID) (Period, error)
}
