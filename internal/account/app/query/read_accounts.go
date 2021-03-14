package query

import "context"

type AccountsReadModel interface {
	AllAccounts(ctx context.Context) ([]Account, error)
	AccountByNumber(ctx context.Context, accountNumber string) (Account, error)
	ValidateExistence(ctx context.Context, accNumbers []string) error
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

func (h ReadAccountsHandler) HandleReadAll(ctx context.Context) ([]Account, error) {
	return h.readModel.AllAccounts(ctx)
}

func (h ReadAccountsHandler) HandleReadByNumber(ctx context.Context, accountNumber string) (Account, error) {
	return h.readModel.AccountByNumber(ctx, accountNumber)
}

func (h ReadAccountsHandler) HandleValidateExistence(ctx context.Context, accNumbers []string) error {
	return h.readModel.ValidateExistence(ctx, accNumbers)
}
