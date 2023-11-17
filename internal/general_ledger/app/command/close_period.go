package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
)

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
	// check all vouchers are posted
	if notPostedVoucherExists, err := h.repo.ExistsVouchersNotPostedInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return fmt.Errorf("failed to check vouchers posted status: %w", err)
	} else if notPostedVoucherExists {
		return commonErrors.NewSlugError("period-close-notAllVouchersPosted")
	}

	// check all profit and loss ledgers have zero ending balance
	if unclearedProfitAndLoss, err := h.repo.ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return fmt.Errorf("failed to check profit and loss ledgers balances: %w", err)
	} else if unclearedProfitAndLoss {
		return commonErrors.NewSlugError("period-close-unclearedProfitAndLoss")
	}

	// check trial balance
	ledgers, err := h.repo.ReadFirstLevelLedgersInPeriod(ctx, cmd.SobId, cmd.PeriodId)
	if err != nil {
		return fmt.Errorf("failed to read 1st level ledgers: %w", err)
	}
	var totalOpeningDebit, totalOpeningCredit,
		totalPeriodDebit, totalPeriodCredit,
		totalEndingDebit, totalEndingCredit decimal.Decimal

	// sum
	for _, l := range ledgers {
		totalOpeningDebit = totalOpeningDebit.Add(l.OpeningDebitBalance())
		totalEndingDebit = totalEndingDebit.Add(l.EndingDebitBalance())

		totalPeriodDebit = totalPeriodDebit.Add(l.PeriodDebit())
		totalPeriodCredit = totalPeriodCredit.Add(l.PeriodCredit())

		totalOpeningCredit = totalOpeningCredit.Add(l.OpeningCreditBalance())
		totalEndingCredit = totalEndingCredit.Add(l.EndingCreditBalance())
	}

	if !totalOpeningDebit.Equal(totalOpeningCredit) {
		return commonErrors.NewSlugError("period-close-openingBalanceUnequal")
	}
	if !totalPeriodDebit.Equal(totalPeriodCredit) {
		return commonErrors.NewSlugError("period-close-periodBalanceUnequal")
	}
	if !totalEndingDebit.Equal(totalEndingCredit) {
		return commonErrors.NewSlugError("period-close-endingBalanceUnequal")
	}

	// update
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.handleUpdate(txCtx, cmd)
	})
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
