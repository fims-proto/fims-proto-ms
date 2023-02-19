package intraprocess

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github/fims-proto/fims-proto-ms/internal/account/app/command"

	"github/fims-proto/fims-proto-ms/internal/account/app"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	"github.com/google/uuid"
)

type AccountInterface struct {
	app *app.Application
}

func NewAccountInterface(app *app.Application) AccountInterface {
	return AccountInterface{app: app}
}

func (i AccountInterface) InitializeAccounts(ctx context.Context, sobId uuid.UUID) error {
	return i.app.Commands.InitialAccounts.Handle(ctx, sobId)
}

func (i AccountInterface) ReadAccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]query.Account, error) {
	return i.app.Queries.AccountsByNumbers.Handle(ctx, sobId, accountNumbers)
}

func (i AccountInterface) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.Account, error) {
	return i.app.Queries.AccountsByIds.Handle(ctx, accountIds)
}

func (i AccountInterface) ReadOrCreatePeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (query.Period, error) {
	p, err := i.app.Queries.PeriodByTime.Handle(ctx, sobId, timePoint)
	if err == nil {
		// found, return
		return p, nil
	} else if err != nil && err.Error() != "period-notFound" {
		// errors others not found
		return query.Period{}, err
	}

	// not found, create

	newPeriodId := uuid.New()
	if err = i.app.Commands.CreateFuturePeriod.Handle(ctx, command.CreateFuturePeriodCmd{
		SobId:     sobId,
		PeriodId:  newPeriodId,
		TimePoint: timePoint,
	}); err != nil {
		return query.Period{}, errors.Wrap(err, "failed to create period when not found")
	}

	return i.app.Queries.PeriodById.Handle(ctx, newPeriodId)
}

func (i AccountInterface) ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]query.Period, error) {
	return i.app.Queries.PeriodsByIds.Handle(ctx, periodIds)
}

func (i AccountInterface) CreatePeriodByNumber(ctx context.Context, cmd command.CreateCurrentPeriodCmd) error {
	return i.app.Commands.CreateCurrentPeriod.Handle(ctx, cmd)
}

func (i AccountInterface) CreateLedgers(ctx context.Context, cmd command.CreateLedgersCmd) error {
	return i.app.Commands.CreateLedgers.Handle(ctx, cmd)
}

func (i AccountInterface) PostAccounts(ctx context.Context, cmd command.PostAccountsCmd) error {
	return i.app.Commands.PostAccounts.Handle(ctx, cmd)
}
