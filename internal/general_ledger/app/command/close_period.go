package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/balance_direction"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
)

type ClosePeriodCmd struct {
	SobId    uuid.UUID
	PeriodId uuid.UUID
}

type ClosePeriodHandler struct {
	repo             domain.Repository
	readModel        query.GeneralLedgerReadModel
	numberingService service.NumberingService
}

func NewClosePeriodHandler(repo domain.Repository, readModel query.GeneralLedgerReadModel, numberingService service.NumberingService) ClosePeriodHandler {
	if repo == nil {
		panic("nil repo")
	}

	if readModel == nil {
		panic("nil read model")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return ClosePeriodHandler{
		repo:             repo,
		readModel:        readModel,
		numberingService: numberingService,
	}
}

func (h ClosePeriodHandler) Handle(ctx context.Context, cmd ClosePeriodCmd) error {
	// check all vouchers are posted
	if notPostedVoucherExists, err := h.readModel.ExistsVouchersNotPostedInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return errors.Wrap(err, "failed to check vouchers posted status")
	} else if notPostedVoucherExists {
		return commonErrors.NewSlugError("period-close-notAllVouchersPosted")
	}

	// check all profit and loss ledgers have zero ending balance
	if unclearedProfitAndLoss, err := h.readModel.ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx, cmd.SobId, cmd.PeriodId); err != nil {
		return errors.Wrap(err, "failed to check profit and loss ledgers balances")
	} else if unclearedProfitAndLoss {
		return commonErrors.NewSlugError("period-close-unclearedProfitAndLoss")
	}

	// check trial balance
	ledgers, err := h.readModel.FirstLevelLedgersInPeriod(ctx, cmd.SobId, cmd.PeriodId)
	if err != nil {
		return errors.Wrap(err, "failed to read 1st level ledgers")
	}
	var totalOpeningDebit, totalOpeningCredit,
		totalPeriodDebit, totalPeriodCredit,
		totalEndingDebit, totalEndingCredit decimal.Decimal

	// sum
	for _, l := range ledgers {
		totalPeriodDebit = totalPeriodDebit.Add(l.PeriodDebit)
		totalPeriodCredit = totalPeriodCredit.Add(l.PeriodCredit)

		if l.Account.BalanceDirection == balance_direction.Debit.String() {
			totalOpeningDebit = totalOpeningDebit.Add(l.OpeningBalance)
			totalEndingDebit = totalEndingDebit.Add(l.EndingBalance)
		} else if l.Account.BalanceDirection == balance_direction.Credit.String() {
			totalOpeningCredit = totalOpeningCredit.Add(l.OpeningBalance)
			totalEndingCredit = totalEndingCredit.Add(l.EndingBalance)
		} else {
			return commonErrors.NewSlugError("period-close-unknownAccountBalanceDirection", l.Account.AccountNumber)
		}
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
	nextPeriodId, err := createPeriodIfNotExists(ctx, createPeriodCmd{
		SobId:      cmd.SobId,
		PeriodId:   uuid.Nil,
		FiscalYear: nextFiscalYear,
		Number:     nextPeriodNumber,
	}, h.repo, h.readModel, h.numberingService)
	if err != nil {
		return errors.Wrap(err, "failed to create next period")
	}

	// update next period to current
	if err = h.repo.UpdatePeriod(ctx, nextPeriodId, func(p *period.Period) (*period.Period, error) {
		if err = p.Start(); err != nil {
			return nil, err
		}
		return p, nil
	}); err != nil {
		return err
	}

	// initialize ledgers for new period
	if err = initializeLedgers(ctx, initializeLedgersCmd{
		SobId:    cmd.SobId,
		PeriodId: nextPeriodId,
	}, h.repo, h.readModel); err != nil {
		return errors.Wrap(err, "failed to initialize ledgers for next period")
	}

	return nil
}
