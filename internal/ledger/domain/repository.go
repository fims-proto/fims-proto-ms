package domain

import "context"

type Repository interface {
	AddLedger(ctx context.Context, l *Ledger) error
	UpdateLedgers(
		ctx context.Context,
		sob string,
		ledgerNumbers []string,
		updateFn func(ledgers []*Ledger) ([]*Ledger, error),
	) error
	Dataload(ctx context.Context, ledgers []*Ledger) error
	Migrate(ctx context.Context) error
}
