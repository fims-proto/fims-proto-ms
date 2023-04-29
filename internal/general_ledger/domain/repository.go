package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
)

type Repository interface {
	InitialAccounts(ctx context.Context, accounts []*account.Account) error

	CreatePeriod(ctx context.Context, period *period.Period) error
	UpdatePeriod(
		ctx context.Context,
		periodId uuid.UUID,
		updateFn func(p *period.Period) (*period.Period, error),
	) error

	CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error
	UpdateLedgersByPeriodAndAccountIds(
		ctx context.Context,
		periodId uuid.UUID,
		accountIds []uuid.UUID,
		updateFn func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error),
	) error

	CreateVoucher(ctx context.Context, d *voucher.Voucher) error
	UpdateVoucher(
		ctx context.Context,
		voucherId uuid.UUID,
		updateFn func(d *voucher.Voucher) (*voucher.Voucher, error),
	) error

	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	Migrate(ctx context.Context) error
}
