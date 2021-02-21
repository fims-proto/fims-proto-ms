package query

import "context"

type accountsReadModel interface {
	AllAccounts(ctx context.Context) ([]Account, error)
	AccountByNumber(ctx context.Context, accountNumber string) (Account, error)
}

type ReadAccountsHandler struct {
	readModel accountsReadModel
}

func NewReadAccountsHandler(readModel accountsReadModel) ReadAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReadAccountsHandler{readModel: readModel}
}

func (h ReadAccountsHandler) HandleReadAll(ctx context.Context) ([]Account, error) {
	return h.readModel.AllAccounts(ctx)
}
