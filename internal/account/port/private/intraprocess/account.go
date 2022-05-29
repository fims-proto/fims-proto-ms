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
	return i.app.Queries.ReadAccounts.HandleReadByAccountNumber(ctx, sobId, accountNumbers)
}

func (i AccountInterface) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	return i.app.Queries.ReadAccounts.HandleReadByIds(ctx, accountIds)
}

func (i AccountInterface) ReadAccountById(ctx context.Context, accountId uuid.UUID) (query.Account, error) {
	return i.app.Queries.ReadAccounts.HandleReadById(ctx, accountId)
}

func (i AccountInterface) ReadAllAccountIdsBySobId(ctx context.Context, sobId uuid.UUID) ([]query.Account, error) {
	// one page that is big enough
	pageRequest, _ := data.NewPageRequest(1, 99999, nil, nil)
	accountsPage, err := i.app.Queries.ReadAccounts.HandleReadAll(ctx, sobId, pageRequest)
	return accountsPage.Content, err
}

func (i AccountInterface) InitializeAccounts(ctx context.Context, sobId uuid.UUID) error {
	return i.app.Commands.LoadAccounts.Handle(ctx, sobId)
}
