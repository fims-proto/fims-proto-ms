package intraprocess

import (
	"context"
	"time"

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

func (i AccountInterface) ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (query.Period, error) {
	return i.app.Queries.PeriodByTime.Handle(ctx, sobId, timePoint)
}

func (i AccountInterface) ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]query.Period, error) {
	return i.app.Queries.PeriodsByIds.Handle(ctx, periodIds)
}

func (i AccountInterface) CreatePeriod(ctx context.Context, cmd command.CreatePeriodCmd) error {
	return i.app.Commands.CreatePeriod.Handle(ctx, cmd)
}

func (i AccountInterface) CreateLedgers(ctx context.Context, cmd command.CreateLedgersCmd) error {
	return i.app.Commands.CreateLedgers.Handle(ctx, cmd)
}

func (i AccountInterface) PostAccounts(ctx context.Context, cmd command.PostAccountsCmd) error {
	return i.app.Commands.PostAccounts.Handle(ctx, cmd)
}
