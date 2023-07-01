package command

import (
	"context"

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

	_, _, err = repo.CreatePeriodIfNotExists(ctx, p)
	return err
}

func createPeriodIfNotExists(
	ctx context.Context,
	cmd createPeriodCmd,
	repo domain.Repository,
	numberingService service.NumberingService,
) (*period.Period, error) {
	p, err := period.New(uuid.New() /*dummy id*/, cmd.SobId, cmd.FiscalYear, cmd.Number, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create period domain model")
	}

	p, created, err := repo.CreatePeriodIfNotExists(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read or create period")
	}

	if created {
		// create numbering configuration for voucher entries in this period
		if err = numberingService.CreateIdentifierConfigurationForVoucher(ctx, p.Id()); err != nil {
			return nil, errors.Wrap(err, "failed to create numbering configuration for period")
		}
	}

	return p, nil
}
