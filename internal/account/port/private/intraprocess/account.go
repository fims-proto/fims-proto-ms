package intraprocess

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

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

func (i AccountInterface) ReadAccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]query.Account, error) {
	return i.app.Queries.ReadAccounts.HandleReadByNumbers(ctx, sobId, accountNumbers)
}

func (i AccountInterface) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	return i.app.Queries.ReadAccounts.HandleReadByIds(ctx, accountIds)
}

func (i AccountInterface) ReadAccountsWithSuperiorsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.Account, error) {
	return i.app.Queries.ReadAccounts.HandleReadWithSuperiorsByIds(ctx, accountIds)
}

func (i AccountInterface) ReadAccountsBySobId(ctx context.Context, sobId uuid.UUID) ([]query.Account, error) {
	accountsPage, err := i.app.Queries.ReadAccounts.HandleReadAll(ctx, sobId, data.Unpaged())
	return accountsPage.Content(), err
}

func (i AccountInterface) InitializeAccounts(ctx context.Context, sobId uuid.UUID) error {
	return i.app.Commands.LoadAccounts.Handle(ctx, sobId)
}
