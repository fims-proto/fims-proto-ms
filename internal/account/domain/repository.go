package domain

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_configuration"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"
)

type Repository interface {
	InitialAccountConfiguration(ctx context.Context, accountConfigurations []*account_configuration.AccountConfiguration) error

	CreatePeriod(ctx context.Context, period *period.Period, txFn func() error) error

	CreateAccounts(ctx context.Context, accounts []*account.Account) error
	UpdateAccountsByPeriodAndIds(
		ctx context.Context,
		periodId uuid.UUID,
		accountIds []uuid.UUID,
		updateFn func(accounts []*account.Account) ([]*account.Account, error),
	) error

	Migrate(ctx context.Context) error
}
