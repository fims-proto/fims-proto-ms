package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateMonthlyClosingJournalCmd struct {
	SobId uuid.UUID
}

type CreateMonthlyClosingJournalHandler struct {
	createJournalHandler CreateJournalHandler
	reviewJournalHandler ReviewJournalHandler
	auditJournalHandler  AuditJournalHandler
	postJournalHandler   PostJournalHandler
	repo                 domain.Repository
}

func NewCreateMonthlyClosingJournalHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	dimensionService service.DimensionService,
	sobService service.SobService,
) CreateMonthlyClosingJournalHandler {
	return CreateMonthlyClosingJournalHandler{
		repo:                 repo,
		createJournalHandler: NewCreateJournalHandler(repo, numberingService, dimensionService, sobService),
		reviewJournalHandler: NewReviewJournalHandler(repo),
		auditJournalHandler:  NewAuditJournalHandler(repo),
		postJournalHandler:   NewPostJournalHandler(repo),
	}
}

// Handle creates, reviews, audits, and posts a monthly closing journal
// that reverses all leaf P&L account balances to zero and transfers
// the net result to the Current Year Profit account (003103).
func (h CreateMonthlyClosingJournalHandler) Handle(ctx context.Context, cmd CreateMonthlyClosingJournalCmd) (uuid.UUID, error) {
	// PRE-CHECKS (before transaction)
	currentPeriod, pnlLedgers, err := h.preCheck(ctx, cmd)
	if err != nil {
		return uuid.Nil, err
	}

	// BUILD JOURNAL LINE COMMANDS
	journalLineCmds, err := h.buildJournal(pnlLedgers)
	if err != nil {
		return uuid.Nil, err
	}

	// SINGLE OUTER TRANSACTION
	journalId := uuid.New()
	if err = h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.process(txCtx, cmd.SobId, journalId, currentPeriod, journalLineCmds)
	}); err != nil {
		return uuid.Nil, err
	}

	return journalId, nil
}

func (h CreateMonthlyClosingJournalHandler) preCheck(ctx context.Context, cmd CreateMonthlyClosingJournalCmd) (*period.Period, []*ledger.Ledger, error) {
	// 1. Read current period
	currentPeriod, err := h.repo.ReadCurrentPeriod(ctx, cmd.SobId)
	if err != nil {
		return nil, nil, commonErrors.NewSlugError("period-notFound")
	}

	// 2. Check if CLOSING journal already exists in this period
	exists, err := h.repo.ExistsClosingJournalInPeriod(ctx, cmd.SobId, currentPeriod.Id(), journal.TypeClosing)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, commonErrors.NewSlugError("journal-closing-alreadyExists")
	}

	// 3. Check all journals are posted
	hasUnposted, err := h.repo.ExistsJournalsNotPostedInPeriod(ctx, cmd.SobId, currentPeriod.Id())
	if err != nil {
		return nil, nil, err
	}
	if hasUnposted {
		return nil, nil, commonErrors.NewSlugError("journal-closing-unpostedJournalsExist")
	}

	// 4. Read P&L ledgers with non-zero balance
	pnlLedgers, err := h.repo.ReadProfitAndLossLedgersHavingBalanceInPeriod(ctx, cmd.SobId, currentPeriod.Id())
	if err != nil {
		return nil, nil, err
	}
	if len(pnlLedgers) == 0 {
		return nil, nil, commonErrors.NewSlugError("journal-closing-noBalanceToClear")
	}
	return currentPeriod, pnlLedgers, nil
}

func (h CreateMonthlyClosingJournalHandler) buildJournal(pnlLedgers []*ledger.Ledger) ([]JournalLineCmd, error) {
	journalLineCmds := make([]JournalLineCmd, 0, len(pnlLedgers)+1)
	var sumPnL decimal.Decimal

	for _, l := range pnlLedgers {
		account := l.Account()
		if account == nil {
			return nil, commonErrors.NewSlugError("account-notFound")
		}

		// Reverse the balance: line.amount = -endingAmount
		journalLineCmds = append(journalLineCmds, JournalLineCmd{
			Id:               uuid.New(),
			RawAccountNumber: account.RawAccountNumber(),
			Text:             "月末结账",
			Amount:           l.EndingAmount().Neg(),
		})

		sumPnL = sumPnL.Add(l.EndingAmount())
	}

	// Add Current Year Profit line: amount = sum of all P&L ending amounts
	journalLineCmds = append(journalLineCmds, JournalLineCmd{
		Id:               uuid.New(),
		RawAccountNumber: "003103", // Current Year Profit (本年利润)
		Text:             "月末结账",
		Amount:           sumPnL,
	})
	return journalLineCmds, nil
}

func (h CreateMonthlyClosingJournalHandler) process(ctx context.Context, sobId, journalId uuid.UUID, currentPeriod *period.Period, journalLineCmds []JournalLineCmd) error {
	txDate := transaction_date.TransactionDate{
		Year:  currentPeriod.FiscalYear(),
		Month: currentPeriod.PeriodNumber(),
		Day:   1,
	}

	// Step 1: Create journal
	if err := h.createJournalHandler.Handle(ctx, CreateJournalCmd{
		JournalId:       journalId,
		SobId:           sobId,
		HeaderText:      "月末结账",
		JournalType:     string(journal.TypeClosing),
		JournalLines:    journalLineCmds,
		Creator:         journal.SystemUser,
		TransactionDate: txDate,
	}); err != nil {
		return err
	}

	// Step 2: Review with SYSTEM user
	if err := h.reviewJournalHandler.Handle(ctx, ReviewJournalCmd{
		JournalId: journalId,
		Reviewer:  journal.SystemUser,
	}); err != nil {
		return err
	}

	// Step 3: Audit with SYSTEM user
	if err := h.auditJournalHandler.Handle(ctx, AuditJournalCmd{
		JournalId: journalId,
		Auditor:   journal.SystemUser,
	}); err != nil {
		return err
	}

	// Step 4: Post with SYSTEM user
	if err := h.postJournalHandler.Handle(ctx, PostJournalCmd{
		JournalId: journalId,
		Poster:    journal.SystemUser,
	}); err != nil {
		return err
	}

	return nil
}
