package command

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github.com/google/uuid"
)

// yearEndRetainedEarningsAccount is the raw account number for 本年利润.
// TODO: make this configurable per SoB in the future.
const yearEndRetainedEarningsAccount = "003103"

type ClosePeriodCmd struct {
	SobId    uuid.UUID
	PeriodId uuid.UUID
}

type ClosePeriodResult struct {
	ReportGenerationFailed bool
}

type ClosePeriodHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
	reportService    service.ReportService // nil when called internally by ClosePeriodsHandler
}

func NewClosePeriodHandler(repo domain.Repository, numberingService service.NumberingService, reportService service.ReportService) ClosePeriodHandler {
	if repo == nil {
		panic("nil repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return ClosePeriodHandler{
		repo:             repo,
		numberingService: numberingService,
		reportService:    reportService,
	}
}

func (h ClosePeriodHandler) Handle(ctx context.Context, cmd ClosePeriodCmd) (ClosePeriodResult, error) {
	// check all journals are posted
	if notPostedJournalExists, err := h.repo.ExistsJournalsNotPostedInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return ClosePeriodResult{}, fmt.Errorf("failed to check journals posted status: %w", err)
	} else if notPostedJournalExists {
		return ClosePeriodResult{}, commonErrors.NewInvalidInputError(commonErrors.SlugPeriodCloseNotAllPosted)
	}

	// check all profit and loss ledgers have zero ending balance
	if unclearedProfitAndLoss, err := h.repo.ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return ClosePeriodResult{}, fmt.Errorf("failed to check profit and loss ledgers balances: %w", err)
	} else if unclearedProfitAndLoss {
		return ClosePeriodResult{}, commonErrors.NewInvalidInputError(commonErrors.SlugPeriodCloseUnclearedPnL)
	}

	// check trial balance
	if err := trialBalance(ctx, h.repo, cmd.SobId, cmd.PeriodId); err != nil {
		return ClosePeriodResult{}, fmt.Errorf("not balance: %w", err)
	}

	// check current-year profit account has zero balance if closing the last period of fiscal year
	if p, err := h.repo.ReadPeriodById(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return ClosePeriodResult{}, fmt.Errorf("failed to read period: %w", err)
	} else if p.PeriodNumber() == 12 {
		if hasBalance, err := h.repo.ExistsLedgerHavingBalanceByRawAccountNumberInPeriod(
			ctx, cmd.SobId, yearEndRetainedEarningsAccount, cmd.PeriodId,
		); err != nil {
			return ClosePeriodResult{}, fmt.Errorf("failed to check year-end account balance: %w", err)
		} else if hasBalance {
			return ClosePeriodResult{}, commonErrors.NewInvalidInputError(commonErrors.SlugPeriodCloseUnclearedProfit)
		}
	}

	// update — period close commits here
	if err := h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.handleUpdate(txCtx, cmd)
	}); err != nil {
		return ClosePeriodResult{}, err
	}

	// generate reports outside the GL transaction; failure surfaces as a warning
	if h.reportService != nil {
		if err := h.reportService.GenerateForPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
			return ClosePeriodResult{ReportGenerationFailed: true}, nil
		}
	}

	return ClosePeriodResult{}, nil
}

func (h ClosePeriodHandler) handleUpdate(ctx context.Context, cmd ClosePeriodCmd) error {
	var nextFiscalYear, nextPeriodNumber int

	// update current period to closed
	if err := h.repo.UpdatePeriod(ctx, cmd.PeriodId, func(p *period.Period) (*period.Period, error) {
		if err := p.Close(); err != nil {
			return nil, err
		}

		// get next period year and number
		nextFiscalYear, nextPeriodNumber = p.NextNumber()

		return p, nil
	}); err != nil {
		return err
	}

	// create next period if it does not exist
	nextPeriod, err := createPeriodIfNotExists(ctx, createPeriodCmd{
		SobId:      cmd.SobId,
		PeriodId:   uuid.Nil,
		FiscalYear: nextFiscalYear,
		Number:     nextPeriodNumber,
	}, h.repo, h.numberingService)
	if err != nil {
		return fmt.Errorf("failed to create next period: %w", err)
	}

	// update next period to current
	if err = h.repo.UpdatePeriod(ctx, nextPeriod.Id(), func(p *period.Period) (*period.Period, error) {
		if err = p.Start(); err != nil {
			return nil, err
		}
		return p, nil
	}); err != nil {
		return err
	}

	// initialize ledgers for new period
	if err = initializeAllLedgers(ctx, h.repo, cmd.SobId); err != nil {
		return fmt.Errorf("failed to initialize ledgers for next period: %w", err)
	}

	return nil
}
