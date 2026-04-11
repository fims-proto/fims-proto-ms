package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github.com/google/uuid"
)

type Repository interface {
	Migrate(ctx context.Context) error
	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	InitialAccounts(ctx context.Context, accounts []*account.Account) error
	CreateAccount(ctx context.Context, account *account.Account) error
	UpdateAccount(
		ctx context.Context,
		accountId uuid.UUID,
		updateFn func(a *account.Account) (*account.Account, error),
	) error
	DeleteAccount(ctx context.Context, accountId uuid.UUID) error
	ReadAllAccounts(ctx context.Context, sobId uuid.UUID) ([]*account.Account, error)
	ReadAccountById(ctx context.Context, accountId uuid.UUID) (*account.Account, error)
	ReadAccountByRawNumber(ctx context.Context, sobId uuid.UUID, rawNumber string) (*account.Account, error)
	ReadAccountsByRawNumbers(ctx context.Context, sobId uuid.UUID, rawNumbers []string) ([]*account.Account, error)
	ReadSuperiorAccountsById(ctx context.Context, accountId uuid.UUID) ([]*account.Account, error)
	ReadAccountsWithSuperiorsByIds(ctx context.Context, sobId uuid.UUID, accountIds []uuid.UUID) ([]*account.Account, error)
	ReadAllSubAccountsWithSuperiors(ctx context.Context, sobId uuid.UUID) ([]*account.Account, error)

	CreatePeriodIfNotExists(ctx context.Context, period *period.Period) (*period.Period, bool, error)
	UpdatePeriod(
		ctx context.Context,
		periodId uuid.UUID,
		updateFn func(p *period.Period) (*period.Period, error),
	) error
	ReadPeriodById(ctx context.Context, sobId, periodId uuid.UUID) (*period.Period, error)
	ReadCurrentPeriod(ctx context.Context, sobId uuid.UUID) (*period.Period, error)
	ReadPreviousPeriod(ctx context.Context, currentPeriodId uuid.UUID) (*period.Period, error)
	ReadFirstPeriod(ctx context.Context, sobId uuid.UUID) (*period.Period, error)

	CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error
	DeleteLedgersByAccountId(ctx context.Context, accountId uuid.UUID) error
	UpdateLedgersByPeriodAndAccountIds(
		ctx context.Context,
		periodId uuid.UUID,
		accountIds []uuid.UUID,
		updateFn func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error),
	) error
	ReadLedgersByPeriod(ctx context.Context, periodId uuid.UUID) ([]*ledger.Ledger, error)
	ReadFirstLevelLedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]*ledger.Ledger, error)
	ReadProfitAndLossLedgersHavingBalanceInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]*ledger.Ledger, error)
	ReadLedgerByRawAccountNumberInPeriod(ctx context.Context, sobId uuid.UUID, rawAccountNumber string, periodId uuid.UUID) (*ledger.Ledger, error)
	ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error)
	ExistsLedgerHavingBalanceByRawAccountNumberInPeriod(ctx context.Context, sobId uuid.UUID, rawAccountNumber string, periodId uuid.UUID) (bool, error)

	CreateJournal(ctx context.Context, j *journal.Journal) error
	ExistsJournalById(ctx context.Context, sobId, journalId uuid.UUID) (bool, error)
	ExistsClosingJournalInPeriod(ctx context.Context, sobId, periodId uuid.UUID, journalType journal.JournalType) (bool, error)
	ExistsJournalLinesByAccountId(ctx context.Context, accountId uuid.UUID) (bool, error)
	ExistsChildAccountsByAccountId(ctx context.Context, accountId uuid.UUID) (bool, error)
	ExistsLedgerWithOpeningBalanceByAccountId(ctx context.Context, accountId uuid.UUID) (bool, error)
	// UpdateJournalHeader updates only the journal header row (status flags, reviewer, auditor, poster, etc).
	// Journal lines are loaded for the callback to read but are NOT deleted or re-saved.
	UpdateJournalHeader(
		ctx context.Context,
		journalId uuid.UUID,
		updateFn func(j *journal.Journal) (*journal.Journal, error),
	) error
	// UpdateEntireJournal replaces the journal header and all its journal lines.
	// Use this when journal lines may be modified; prefer UpdateJournalHeader for header-only changes.
	UpdateEntireJournal(
		ctx context.Context,
		journalId uuid.UUID,
		updateFn func(j *journal.Journal) (*journal.Journal, error),
	) error
	ExistsJournalsNotPostedInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error)

	// ReadJournalById reads a journal (header + lines + period) by its primary key.
	// Returns commonErrors.ErrRecordNotFound() if no row exists.
	ReadJournalById(ctx context.Context, journalId uuid.UUID) (*journal.Journal, error)

	// DeleteJournalById hard-deletes the journal header, all journal_lines rows,
	// and all journal_line_dimension_options rows.
	// Callers must reverse ledgers before calling this method.
	DeleteJournalById(ctx context.Context, journalId uuid.UUID) error
}
