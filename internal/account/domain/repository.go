package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/account/domain/ledger"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"
)

type Repository interface {
	InitialAccounts(ctx context.Context, accounts []*account.Account) error

	CreatePeriod(ctx context.Context, period *period.Period, txFn func() error) error

	CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error
	UpdateLedgersByPeriodAndAccountIds(
		ctx context.Context,
		periodId uuid.UUID,
		accountIds []uuid.UUID,
		updateFn func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error),
	) error

	Migrate(ctx context.Context) error
}
