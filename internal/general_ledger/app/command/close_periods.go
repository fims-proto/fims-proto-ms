package command

import (
	"context"
	"fmt"
	"time"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type ClosePeriodsCmd struct {
	SobId       uuid.UUID
	TargetYear  int
	TargetMonth int
}

type ClosePeriodsHandler struct {
	repo                  domain.Repository
	closePeriodHandler    ClosePeriodHandler
	monthlyClosingHandler CreateMonthlyClosingJournalHandler
	yearEndClosingHandler CreateYearEndClosingJournalHandler
}

func NewClosePeriodsHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	dimensionService service.DimensionService,
	sobService service.SobService,
) ClosePeriodsHandler {
	if repo == nil {
		panic("nil repo")
	}
	return ClosePeriodsHandler{
		repo:                  repo,
		closePeriodHandler:    NewClosePeriodHandler(repo, numberingService),
		monthlyClosingHandler: NewCreateMonthlyClosingJournalHandler(repo, numberingService, dimensionService, sobService),
		yearEndClosingHandler: NewCreateYearEndClosingJournalHandler(repo, numberingService, dimensionService, sobService),
	}
}

func (h ClosePeriodsHandler) Handle(ctx context.Context, cmd ClosePeriodsCmd) error {
	current, err := h.repo.ReadCurrentPeriod(ctx, cmd.SobId)
	if err != nil {
		return commonErrors.NewInvalidInputError(commonErrors.SlugPeriodNotFound)
	}

	if cmd.TargetYear < current.FiscalYear() ||
		(cmd.TargetYear == current.FiscalYear() && cmd.TargetMonth < current.PeriodNumber()) {
		return commonErrors.NewInvalidInputError(commonErrors.SlugPeriodBatchCloseTargetInPast)
	}

	sequence := buildPeriodSequence(current.FiscalYear(), current.PeriodNumber(), cmd.TargetYear, cmd.TargetMonth)
	if len(sequence) > 12 {
		return commonErrors.NewInvalidInputError(commonErrors.SlugPeriodBatchCloseTooManyPeriods)
	}

	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		for _, ref := range sequence {
			if err := h.closeSinglePeriod(txCtx, cmd.SobId, ref.year, ref.month); err != nil {
				return err
			}
		}
		return nil
	})
}

func (h ClosePeriodsHandler) closeSinglePeriod(ctx context.Context, sobId uuid.UUID, expectedYear, expectedMonth int) error {
	current, err := h.repo.ReadCurrentPeriod(ctx, sobId)
	if err != nil {
		return fmt.Errorf("failed to read current period: %w", err)
	}
	if current.FiscalYear() != expectedYear || current.PeriodNumber() != expectedMonth {
		return fmt.Errorf("unexpected current period: got %d-%02d, expected %d-%02d",
			current.FiscalYear(), current.PeriodNumber(), expectedYear, expectedMonth)
	}

	// Create monthly closing journal (skip if P&L has no balance).
	pnlLedgers, err := h.repo.ReadProfitAndLossLedgersHavingBalanceInPeriod(ctx, sobId, current.Id())
	if err != nil {
		return fmt.Errorf("failed to check P&L balance: %w", err)
	}
	if len(pnlLedgers) > 0 {
		if _, err = h.monthlyClosingHandler.Handle(ctx, CreateMonthlyClosingJournalCmd{SobId: sobId}); err != nil {
			return fmt.Errorf("failed to create monthly closing journal for %d-%02d: %w", expectedYear, expectedMonth, err)
		}
	}

	// Create year-end closing journal for period 12 (skip if CYP has no balance).
	if current.PeriodNumber() == 12 {
		cypLedger, err := h.repo.ReadLedgerByRawAccountNumberInPeriod(ctx, sobId, yearEndRetainedEarningsAccount, current.Id())
		if err != nil {
			return fmt.Errorf("failed to read CYP ledger: %w", err)
		}
		if cypLedger != nil && !cypLedger.EndingAmount().IsZero() {
			if _, err = h.yearEndClosingHandler.Handle(ctx, CreateYearEndClosingJournalCmd{SobId: sobId}); err != nil {
				return fmt.Errorf("failed to create year-end closing journal for %d-%02d: %w", expectedYear, expectedMonth, err)
			}
		}
	}

	// Close the period (re-validates all checks — they should pass after auto-journals).
	return h.closePeriodHandler.Handle(ctx, ClosePeriodCmd{
		SobId:    sobId,
		PeriodId: current.Id(),
	})
}

type periodRef struct {
	year  int
	month int
}

func buildPeriodSequence(fromYear, fromMonth, toYear, toMonth int) []periodRef {
	var seq []periodRef
	year, month := fromYear, fromMonth
	for year <= toYear && (year != toYear || month <= toMonth) {
		seq = append(seq, periodRef{year, month})
		t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
		year, month = t.Year(), int(t.Month())
	}
	return seq
}
