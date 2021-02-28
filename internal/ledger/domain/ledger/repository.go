package ledger

import "context"

type Repository interface {
	UpdateLedgers(
		ctx context.Context,
		ledgerNumbers []string,
		updateFn func(ledgers []*Ledger) ([]*Ledger, error),
	) error
}
