package query

import "context"

type AllAccountsReadModel interface {
	AllAccounts(ctx context.Context) ([]Account, error)
}

type AllAccountsHandler struct {
	readModel AllAccountsReadModel
}

func NewAllAccountsHandler(readModel AllAccountsReadModel) AllAccountsHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return AllAccountsHandler{readModel: readModel}
}

func (h AllAccountsHandler) Handle(ctx context.Context) ([]Account, error) {
	return h.readModel.AllAccounts(ctx)
}
