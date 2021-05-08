package query

import "context"

type CountersReadModel interface {
	CounterByBusinessObject(ctx context.Context, businessObject string) (Counter, error)
}

type ReadCountersHandler struct {
	readModel CountersReadModel
}

func NewReadCountersHandler(readModel CountersReadModel) ReadCountersHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReadCountersHandler{readModel: readModel}
}

func (h ReadCountersHandler) HandleByBusinessObject(ctx context.Context, businessObject string) (Counter, error) {
	return h.readModel.CounterByBusinessObject(ctx, businessObject)
}
