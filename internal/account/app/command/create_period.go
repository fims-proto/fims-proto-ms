package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"

	"github/fims-proto/fims-proto-ms/internal/account/app/service"

	"github/fims-proto/fims-proto-ms/internal/account/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// create first period in the SoB
// create period, checking number, using ending time as opening time

type CreatePeriodCmd struct {
	SobId      uuid.UUID
	PeriodId   uuid.UUID
	FiscalYear int
	Number     int
}

type CreatePeriodHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
	readModel        query.AccountReadModel
}

func NewCreatePeriodHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	readModel query.AccountReadModel,
) CreatePeriodHandler {
	if repo == nil {
		panic("nil account repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	if readModel == nil {
		panic("nil read model")
	}

	return CreatePeriodHandler{
		repo:             repo,
		numberingService: numberingService,
		readModel:        readModel,
	}
}

func (h CreatePeriodHandler) Handle(ctx context.Context, cmd CreatePeriodCmd) error {
	p, err := period.New(cmd.PeriodId, cmd.SobId, cmd.FiscalYear, cmd.Number, false)
	if err != nil {
		return errors.Wrap(err, "failed to create period domain model")
	}

	return h.repo.CreatePeriod(ctx, p, func() error {
		// create numbering configuration for voucher entries in this period
		if err = h.numberingService.InitializeIdentifierConfigurationForVoucher(ctx, cmd.PeriodId); err != nil {
			return errors.Wrap(err, "failed to create numbering configuration for period")
		}

		return nil
	})
}
