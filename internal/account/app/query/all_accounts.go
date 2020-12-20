package query

import "context"

type AllAccountsHandler struct {
	readModel AllAccountsReadModel
}

func NewAllAccountsHandler(readModel AllAccountsReadModel) AllAccountsHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return AllAccountsHandler{readModel: readModel}
}

func (h AllAccountsHandler) handle(ctx context.Context) ([]Account, error) {
	return h.readModel.AllAccounts(ctx)
}

type AllAccountsReadModel interface {
	AllAccounts(ctx context.Context) ([]Account, error)
}
