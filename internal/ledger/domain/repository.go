package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreatePeriod(ctx context.Context, period *Period) error

	CreateLedgers(ctx context.Context, ledgers []*Ledger) error
	UpdateLedgersByPeriodAndAccounts(ctx context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(ledgers []*Ledger) ([]*Ledger, error)) error

	Migrate(ctx context.Context) error
}
