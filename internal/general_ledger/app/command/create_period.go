package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
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
		return fmt.Errorf("failed to create period: %w", err)
	}

	// create numbering configuration for voucher entries in this period
	if err = numberingService.CreateIdentifierConfigurationForVoucher(ctx, cmd.PeriodId); err != nil {
		return fmt.Errorf("failed to create period: %w", err)
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
		return nil, fmt.Errorf("failed to create period: %w", err)
	}

	p, created, err := repo.CreatePeriodIfNotExists(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to create period: %w", err)
	}

	if created {
		// create numbering configuration for voucher entries in this period
		if err = numberingService.CreateIdentifierConfigurationForVoucher(ctx, p.Id()); err != nil {
			return nil, fmt.Errorf("failed to create period: %w", err)
		}
	}

	return p, nil
}
