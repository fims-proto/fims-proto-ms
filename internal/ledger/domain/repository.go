package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreatePeriod(ctx context.Context, period *Period) error
	UpdatePeriod(ctx context.Context, id uuid.UUID, updateFn func(period *Period) (*Period, error)) error
	CreateLedgers(ctx context.Context, ledgers []*Ledger) error
	UpdateLedgers(ctx context.Context, ids []uuid.UUID, updateFn func(ledgers []*Ledger) ([]*Ledger, error)) error
	UpdatePeriodLedgers(ctx context.Context, periodId uuid.UUID, updateFn func(ledgers []*Ledger) ([]*Ledger, error)) error
	CreateLedgerLogs(ctx context.Context, logs []*LedgerLog) error
	Migrate(ctx context.Context) error
}
