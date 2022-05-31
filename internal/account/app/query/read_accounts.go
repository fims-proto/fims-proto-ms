package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AccountsReadModel interface {
	ReadAllAccounts(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Account], error)
	ReadById(ctx context.Context, accountId uuid.UUID) (Account, error)
	ReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]Account, error)
	ReadByAccountNumber(ctx context.Context, sobId uuid.UUID, accountNumber string) (Account, error)
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
	return h.readModel.ReadAllAccounts(ctx, sobId, pageable)
}

func (h ReadAccountsHandler) HandleReadById(ctx context.Context, accountId uuid.UUID) (Account, error) {
	return h.readModel.ReadById(ctx, accountId)
}

func (h ReadAccountsHandler) HandleReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]Account, error) {
	return h.readModel.ReadByIds(ctx, accountIds)
}

func (h ReadAccountsHandler) HandleReadByAccountNumber(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]Account, error) {
	accounts := make(map[string]Account)

	for _, accountNumber := range accountNumbers {
		account, err := h.readModel.ReadByAccountNumber(ctx, sobId, accountNumber)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		accounts[accountNumber] = account
	}
	return accounts, nil
}
