package query

import (
	"context"

	"github.com/pkg/errors"
)

type AccountsReadModel interface {
	ReadAllAccounts(ctx context.Context, sob string) ([]Account, error)
	ReadByNumber(ctx context.Context, sob, accountNumber string) (Account, error)
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

func (h ReadAccountsHandler) HandleReadAll(ctx context.Context, sob string) ([]Account, error) {
	return h.readModel.ReadAllAccounts(ctx, sob)
}

func (h ReadAccountsHandler) HandleReadByNumber(ctx context.Context, sob, accountNumber string) (Account, error) {
	return h.readModel.ReadByNumber(ctx, sob, accountNumber)
}

func (h ReadAccountsHandler) HandleValidateExistence(ctx context.Context, sob string, accNumbers []string) error {
	for _, accountNumber := range accNumbers {
		if _, err := h.readModel.ReadByNumber(ctx, sob, accountNumber); err != nil {
			return errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
	}
	return nil
}
