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
)

type CreateYearEndClosingJournalCmd struct {
	SobId uuid.UUID
}

type CreateYearEndClosingJournalHandler struct {
	createJournalHandler CreateJournalHandler
	reviewJournalHandler ReviewJournalHandler
	auditJournalHandler  AuditJournalHandler
	postJournalHandler   PostJournalHandler
	repo                 domain.Repository
}

func NewCreateYearEndClosingJournalHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	dimensionService service.DimensionService,
	sobService service.SobService,
) CreateYearEndClosingJournalHandler {
	return CreateYearEndClosingJournalHandler{
		repo:                 repo,
		createJournalHandler: NewCreateJournalHandler(repo, numberingService, dimensionService, sobService),
		reviewJournalHandler: NewReviewJournalHandler(repo),
		auditJournalHandler:  NewAuditJournalHandler(repo),
		postJournalHandler:   NewPostJournalHandler(repo),
	}
}

// Handle creates, reviews, audits, and posts a year-end closing journal
// that transfers the Current Year Profit balance to Retained Earnings.
// Only callable in period 12 (year-end) and only after monthly closing is done.
func (h CreateYearEndClosingJournalHandler) Handle(ctx context.Context, cmd CreateYearEndClosingJournalCmd) (uuid.UUID, error) {
	// PRE-CHECKS (before transaction)
	currentPeriod, cypLedger, err := h.preCheck(ctx, cmd)
	if err != nil {
		return uuid.Nil, err
	}

	// BUILD JOURNAL LINE COMMANDS
	journalLineCmds := h.buildJournal(cypLedger)

	// SINGLE OUTER TRANSACTION
	journalId := uuid.New()
	if err = h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.process(txCtx, cmd.SobId, journalId, currentPeriod, journalLineCmds)
	}); err != nil {
		return uuid.Nil, err
	}

	return journalId, nil
}

func (h CreateYearEndClosingJournalHandler) preCheck(ctx context.Context, cmd CreateYearEndClosingJournalCmd) (*period.Period, *ledger.Ledger, error) {
	// 1. Read current period
	currentPeriod, err := h.repo.ReadCurrentPeriod(ctx, cmd.SobId)
	if err != nil {
		return nil, nil, commonErrors.NewSlugError("period-notFound")
	}

	// 2. Validate period is month 12 (year-end)
	if currentPeriod.PeriodNumber() != 12 {
		return nil, nil, commonErrors.NewSlugError("journal-yearEndClosing-notYearEndPeriod")
	}

	// 3. Check if YEARLY_CLOSING journal already exists
	exists, err := h.repo.ExistsClosingJournalInPeriod(ctx, cmd.SobId, currentPeriod.Id(), journal.TypeYearlyClosing)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, commonErrors.NewSlugError("journal-yearEndClosing-alreadyExists")
	}

	// 4. Check monthly closing was done (CLOSING journal must exist)
	monthlyDone, err := h.repo.ExistsClosingJournalInPeriod(ctx, cmd.SobId, currentPeriod.Id(), journal.TypeClosing)
	if err != nil {
		return nil, nil, err
	}
	if !monthlyDone {
		return nil, nil, commonErrors.NewSlugError("journal-yearEndClosing-monthlyClosingRequired")
	}

	// 5. Read Current Year Profit ledger
	cypLedger, err := h.repo.ReadLedgerByRawAccountNumberInPeriod(ctx, cmd.SobId, "003103", currentPeriod.Id())
	if err != nil {
		return nil, nil, err
	}
	if cypLedger == nil {
		return nil, nil, commonErrors.NewSlugError("account-notFound")
	}

	// Check if CYP balance is zero
	if cypLedger.EndingAmount().IsZero() {
		return nil, nil, commonErrors.NewSlugError("journal-yearEndClosing-noBalanceToClear")
	}
	return currentPeriod, cypLedger, nil
}

func (h CreateYearEndClosingJournalHandler) buildJournal(cypLedger *ledger.Ledger) []JournalLineCmd {
	amount := cypLedger.EndingAmount()
	journalLineCmds := []JournalLineCmd{
		{
			Id:               uuid.New(),
			RawAccountNumber: "003103", // Current Year Profit (本年利润)
			Text:             "年末结账",
			Amount:           amount.Neg(), // Reverse the balance
		},
		{
			Id:               uuid.New(),
			RawAccountNumber: "003104000002", // Retained Earnings (未分配利润)
			Text:             "年末结账",
			Amount:           amount, // Transfer the amount
		},
	}
	return journalLineCmds
}

func (h CreateYearEndClosingJournalHandler) process(ctx context.Context, sobId, journalId uuid.UUID, currentPeriod *period.Period, journalLineCmds []JournalLineCmd) error {
	txDate := transaction_date.TransactionDate{
		Year:  currentPeriod.FiscalYear(),
		Month: 12,
		Day:   31,
	}

	// Step 1: Create journal
	if err := h.createJournalHandler.Handle(ctx, CreateJournalCmd{
		JournalId:       journalId,
		SobId:           sobId,
		HeaderText:      "年末结账",
		JournalType:     string(journal.TypeYearlyClosing),
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
