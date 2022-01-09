package counter

import (
	"context"
	counterPort "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	counterInterface counterPort.CounterInterface
}

func NewIntraProcessAdapter(counterInterface counterPort.CounterInterface) IntraProcessAdapter {
	return IntraProcessAdapter{counterInterface: counterInterface}
}

func (s IntraProcessAdapter) GetNextIdentifier(ctx context.Context, businessObjects ...string) (string, error) {
	return s.counterInterface.Next(ctx, "", businessObjects)
}
