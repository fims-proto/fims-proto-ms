package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"

	"github/fims-proto/fims-proto-ms/internal/account/app/service"

	"github/fims-proto/fims-proto-ms/internal/account/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateFuturePeriodCmd struct {
	SobId     uuid.UUID
	PeriodId  uuid.UUID
	TimePoint time.Time
}

type CreateFuturePeriodHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
	readModel        query.AccountReadModel
}

func NewCreateFuturePeriodHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	readModel query.AccountReadModel,
) CreateFuturePeriodHandler {
	if repo == nil {
		panic("nil account repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	if readModel == nil {
		panic("nil read model")
	}

	return CreateFuturePeriodHandler{
		repo:             repo,
		numberingService: numberingService,
		readModel:        readModel,
	}
}

func (h CreateFuturePeriodHandler) Handle(ctx context.Context, cmd CreateFuturePeriodCmd) error {
	p, err := period.NewFuture(cmd.PeriodId, cmd.SobId, cmd.TimePoint)
	if err != nil {
		return errors.Wrap(err, "failed to create period domain model")
	}

	return h.repo.CreatePeriod(ctx, p, func() error {
		// create numbering configuration for voucher entries in this period
		if err = h.numberingService.CreateIdentifierConfigurationForVoucher(ctx, cmd.PeriodId); err != nil {
			return errors.Wrap(err, "failed to create numbering configuration for period")
		}

		return nil
	})
}
