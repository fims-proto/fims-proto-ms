package counter

import (
	"context"

	"github.com/google/uuid"
	counterPort "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	counterInterface counterPort.CounterInterface
}

func NewIntraProcessAdapter(counterInterface counterPort.CounterInterface) IntraProcessAdapter {
	return IntraProcessAdapter{counterInterface: counterInterface}
}

func (i IntraProcessAdapter) InitializeCounters(ctx context.Context, sobId uuid.UUID) error {
	return i.counterInterface.InitializeCounters(ctx, sobId)
}
