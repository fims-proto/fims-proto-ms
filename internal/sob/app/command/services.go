package command

import (
	"context"

	"github.com/google/uuid"
)

type AccountService interface {
	InitializeAccounts(ctx context.Context, sobId uuid.UUID) error
}

type LedgerService interface {
	InitializeFirstPeriod(ctx context.Context, sobId uuid.UUID, financialYear, number int) error
}

type CounterService interface {
	InitializeCounters(ctx context.Context, sobId uuid.UUID) error
}
