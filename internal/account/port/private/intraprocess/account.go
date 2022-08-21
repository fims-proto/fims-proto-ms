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

func (i AccountInterface) InitializeAccountConfigurations(ctx context.Context, sobId uuid.UUID) error {
	return i.app.Commands.InitialAccountConfigurations.Handle(ctx, sobId)
}

func (i AccountInterface) ReadAccountConfigurationsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]query.AccountConfiguration, error) {
	return i.app.Queries.AccountConfigurationsByNumbers.Handle(ctx, sobId, accountNumbers)
}

func (i AccountInterface) ReadAccountConfigurationsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.AccountConfiguration, error) {
	return i.app.Queries.AccountConfigurationsByIds.Handle(ctx, accountIds)
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

func (i AccountInterface) PostAccounts(ctx context.Context, cmd command.PostAccountsCmd) error {
	return i.app.Commands.PostAccounts.Handle(ctx, cmd)
}
