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

type ClosePeriodHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
}

func NewClosePeriodHandler(repo domain.Repository, numberingService service.NumberingService) ClosePeriodHandler {
	if repo == nil {
		panic("nil repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return ClosePeriodHandler{
		repo:             repo,
		numberingService: numberingService,
	}
}

func (h ClosePeriodHandler) Handle(ctx context.Context, cmd ClosePeriodCmd) error {
	// check all journals are posted
	if notPostedJournalExists, err := h.repo.ExistsJournalsNotPostedInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return fmt.Errorf("failed to check journals posted status: %w", err)
	} else if notPostedJournalExists {
		return commonErrors.NewSlugError("period-close-notAllJournalsPosted")
	}

	// check all profit and loss ledgers have zero ending balance
	if unclearedProfitAndLoss, err := h.repo.ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return fmt.Errorf("failed to check profit and loss ledgers balances: %w", err)
	} else if unclearedProfitAndLoss {
		return commonErrors.NewSlugError("period-close-unclearedProfitAndLoss")
	}

	// check trial balance
	if err := trialBalance(ctx, h.repo, cmd.SobId, cmd.PeriodId); err != nil {
		return fmt.Errorf("not balance: %w", err)
	}

	// update
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.handleUpdate(txCtx, cmd)
	})
}

func (h ClosePeriodHandler) handleUpdate(ctx context.Context, cmd ClosePeriodCmd) error {
	var nextFiscalYear, nextPeriodNumber int
	var closingPeriodNumber int

	// update current period to closed
	if err := h.repo.UpdatePeriod(ctx, cmd.PeriodId, func(p *period.Period) (*period.Period, error) {
		if err := p.Close(); err != nil {
			return nil, err
		}

		// get next period year and number
		nextFiscalYear, nextPeriodNumber = p.NextNumber()
		closingPeriodNumber = p.PeriodNumber()

		return p, nil
	}); err != nil {
		return err
	}

	// if closing the last period of a fiscal year, check 本年利润 has zero balance
	if closingPeriodNumber == 12 {
		if hasBalance, err := h.repo.ExistsLedgerHavingBalanceByRawAccountNumberInPeriod(
			ctx, cmd.SobId, yearEndRetainedEarningsAccount, cmd.PeriodId,
		); err != nil {
			return fmt.Errorf("failed to check year-end account balance: %w", err)
		} else if hasBalance {
			return commonErrors.NewSlugError("period-close-unclearedYearEndAccount")
		}
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
