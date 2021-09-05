package query

import "context"

type LedgersReadModel interface {
	ReadAllLedgers(ctx context.Context, sob string) ([]Ledger, error)
}

type ReadLedgersHandler struct {
	readModel LedgersReadModel
}

func NewReadLedgersHandler(readModel LedgersReadModel) ReadLedgersHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReadLedgersHandler{readModel: readModel}
}

func (h ReadLedgersHandler) HandleReadAll(ctx context.Context, sob string) ([]Ledger, error) {
	return h.readModel.ReadAllLedgers(ctx, sob)
}
