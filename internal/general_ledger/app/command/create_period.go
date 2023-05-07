package command

import (
	"context"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type createPeriodCmd struct {
	SobId      uuid.UUID
	PeriodId   uuid.UUID
	FiscalYear int
	Number     int
}

func createPeriod(
	ctx context.Context,
	cmd createPeriodCmd,
	repo domain.Repository,
	numberingService service.NumberingService,
) error {
	p, err := period.New(cmd.PeriodId, cmd.SobId, cmd.FiscalYear, cmd.Number, true)
	if err != nil {
		return errors.Wrap(err, "failed to create period domain model")
	}

	// create numbering configuration for voucher entries in this period
	if err = numberingService.CreateIdentifierConfigurationForVoucher(ctx, cmd.PeriodId); err != nil {
		return errors.Wrap(err, "failed to create numbering configuration for period")
	}

	return repo.CreatePeriod(ctx, p)
}

func createPeriodIfNotExists(
	ctx context.Context,
	cmd createPeriodCmd,
	repo domain.Repository,
	readModel query.GeneralLedgerReadModel,
	numberingService service.NumberingService,
) (uuid.UUID, error) {
	existedPeriod, err := readModel.PeriodByFiscalYearAndNumber(ctx, cmd.SobId, cmd.FiscalYear, cmd.Number)
	if err == nil {
		// period exists
		return existedPeriod.Id, nil
	}
	if _, ok := err.(commonErrors.ObjectNotFoundErr); !ok {
		return uuid.Nil, errors.Wrap(err, "failed to read period by transaction time")
	}

	// create period
	newPeriodId := uuid.New()

	p, err := period.New(newPeriodId, cmd.SobId, cmd.FiscalYear, cmd.Number, false)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to create period domain model")
	}

	// create numbering configuration for voucher entries in this period
	if err = numberingService.CreateIdentifierConfigurationForVoucher(ctx, newPeriodId); err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to create numbering configuration for period")
	}

	if err = repo.CreatePeriod(ctx, p); err != nil {
		return uuid.Nil, err
	} else {
		return newPeriodId, nil
	}
}
