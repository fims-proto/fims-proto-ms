package account

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	accountPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	accountInterface accountPort.AccountInterface
}

func NewIntraProcessAdapter(accountInterface accountPort.AccountInterface) IntraProcessAdapter {
	return IntraProcessAdapter{accountInterface: accountInterface}
}

func (s IntraProcessAdapter) ValidateExistenceAndGetId(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error) {
	accounts, err := s.accountInterface.ReadAccountsByNumbers(ctx, sobId, accountNumbers)
	if err != nil {
		return nil, errors.Wrap(err, "validate existence failed")
	}
	accountIds := make(map[string]uuid.UUID)
	for accountNumber, account := range accounts {
		accountIds[accountNumber] = account.Id
	}
	return accountIds, nil
}

func (s IntraProcessAdapter) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	return s.accountInterface.ReadAccountsByIds(ctx, accountIds)
}
