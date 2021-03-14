package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app"

	"github.com/pkg/errors"
)

type AccountInterface struct {
	app *app.Application
}

func NewAccountInterface(app *app.Application) AccountInterface {
	return AccountInterface{app: app}
}

func (i AccountInterface) ValidateExistence(ctx context.Context, accNumbers []string) error {
	return i.app.Queries.ReadAccounts.HandleValidateExistence(ctx, accNumbers)
}

func (i AccountInterface) ReadSuperiorNumbers(ctx context.Context, accNumber string) ([]string, error) {
	acc, err := i.app.Queries.ReadAccounts.HandleReadByNumber(ctx, accNumber)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read account by number %s", accNumber)
	}

	// read only numbers
	accNums := []string{}
	account := &acc
	for {
		if account == nil {
			break
		}
		accNums = append(accNums, account.Number)
		account = account.SuperiorAccount
	}
	return accNums, nil
}
