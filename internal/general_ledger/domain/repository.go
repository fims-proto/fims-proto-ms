package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
)

type Repository interface {
	Migrate(ctx context.Context) error
	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	InitialAccounts(ctx context.Context, accounts []*account.Account) error
	UpdateAccount(
		ctx context.Context,
		accountId uuid.UUID,
		updateFn func(a *account.Account) (*account.Account, error),
	) error
	ReadAllAccounts(ctx context.Context, sobId uuid.UUID) ([]*account.Account, error)
	ReadAccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]*account.Account, error)
	ReadSuperiorAccountsById(ctx context.Context, accountId uuid.UUID) ([]*account.Account, error)
	ReadAccountsWithSuperiorsByIds(ctx context.Context, sobId uuid.UUID, accountIds []uuid.UUID) ([]*account.Account, error)
	ReadAllSubAccountsWithSuperiors(ctx context.Context, sobId uuid.UUID) ([]*account.Account, error)

	CreatePeriodIfNotExists(ctx context.Context, period *period.Period) (*period.Period, bool, error)
	UpdatePeriod(
		ctx context.Context,
		periodId uuid.UUID,
		updateFn func(p *period.Period) (*period.Period, error),
	) error
	ReadCurrentPeriod(ctx context.Context, sobId uuid.UUID) (*period.Period, error)
	ReadPreviousPeriod(ctx context.Context, currentPeriodId uuid.UUID) (*period.Period, error)
	ReadFirstPeriod(ctx context.Context, sobId uuid.UUID) (*period.Period, error)

	CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error
	UpdateLedgersByPeriodAndAccountIds(
		ctx context.Context,
		periodId uuid.UUID,
		accountIds []uuid.UUID,
		updateFn func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error),
	) error
	ReadLedgersByPeriod(ctx context.Context, periodId uuid.UUID) ([]*ledger.Ledger, error)
	ReadFirstLevelLedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]*ledger.Ledger, error)
	ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error)

	CreateVoucher(ctx context.Context, v *voucher.Voucher) error
	UpdateVoucher(
		ctx context.Context,
		voucherId uuid.UUID,
		updateFn func(v *voucher.Voucher) (*voucher.Voucher, error),
	) error
	ExistsVouchersNotPostedInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error)

	CreateAuxiliaryCategories(ctx context.Context, categories []*auxiliary_category.AuxiliaryCategory) error
	ReadAuxiliaryCategoryByKey(ctx context.Context, key string) (*auxiliary_category.AuxiliaryCategory, error)

	CreateAuxiliaryAccounts(ctx context.Context, accounts []*auxiliary_account.AuxiliaryAccount) error
	ReadAuxiliaryAccountsByPairs(ctx context.Context, sobId uuid.UUID, pairs []auxiliary_account.AuxiliaryPair) ([]*auxiliary_account.AuxiliaryAccount, error)
	ReadAllAuxiliaryAccounts(ctx context.Context, sobId uuid.UUID) ([]*auxiliary_account.AuxiliaryAccount, error)

	CreateAuxiliaryLedgers(ctx context.Context, ledgers []*auxiliary_ledger.AuxiliaryLedger) error
	UpdateAuxiliaryLedgersByPeriodAndAccountIds(
		ctx context.Context,
		periodId uuid.UUID,
		auxiliaryAccountIds []uuid.UUID,
		updateFn func(auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger) ([]*auxiliary_ledger.AuxiliaryLedger, error),
	) error
	ReadAuxiliaryLedgersByPeriod(ctx context.Context, periodId uuid.UUID) ([]*auxiliary_ledger.AuxiliaryLedger, error)
}
