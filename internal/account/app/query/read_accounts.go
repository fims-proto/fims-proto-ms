package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AccountsReadModel interface {
	ReadAccounts(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Account], error)
	ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]Account, error)
	ReadAccountsWithSuperiorsByIds(ctx context.Context, accountIds []uuid.UUID) ([]Account, error)
	ReadAccountByNumber(ctx context.Context, sobId uuid.UUID, accountNumber string) (Account, error)
}

type ReadAccountsHandler struct {
	readModel AccountsReadModel
}

func NewReadAccountsHandler(readModel AccountsReadModel) ReadAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReadAccountsHandler{readModel: readModel}
}

func (h ReadAccountsHandler) HandleReadAll(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Account], error) {
	return h.readModel.ReadAccounts(ctx, sobId, pageable)
}

func (h ReadAccountsHandler) HandleReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]Account, error) {
	return h.readModel.ReadAccountsByIds(ctx, accountIds)
}

func (h ReadAccountsHandler) HandleReadWithSuperiorsByIds(ctx context.Context, accountIds []uuid.UUID) ([]Account, error) {
	return h.readModel.ReadAccountsWithSuperiorsByIds(ctx, accountIds)
}

func (h ReadAccountsHandler) HandleReadByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]Account, error) {
	accounts := make(map[string]Account)

	for _, accountNumber := range accountNumbers {
		account, err := h.readModel.ReadAccountByNumber(ctx, sobId, accountNumber)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		accounts[accountNumber] = account
	}
	return accounts, nil
}
