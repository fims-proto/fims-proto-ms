package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAccountingPeriod(ctx context.Context, period *AccountingPeriod) error
	UpdateAccountingPeriod(ctx context.Context, id uuid.UUID, updateFn func(period *AccountingPeriod) (*AccountingPeriod, error)) error
	CreateLedgers(ctx context.Context, ledgers []*Ledger) error
	UpdateLedgers(ctx context.Context, ids []uuid.UUID, updateFn func(ledgers []*Ledger) ([]*Ledger, error)) error
	CreateLedgerLogs(ctx context.Context, logs []*LedgerLog) error
	Migrate(ctx context.Context) error
}
