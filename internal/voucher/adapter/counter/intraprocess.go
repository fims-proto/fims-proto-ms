package counter

import (
	"context"
	counterport "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"
)

type IntraprocessAdapter struct {
	cntrInterface counterport.CounterInterface
}

func NewIntraprocessAdapter(cntrInterface counterport.CounterInterface) IntraprocessAdapter {
	return IntraprocessAdapter{cntrInterface: cntrInterface}
}

func (s IntraprocessAdapter) GetNextIdentifier(ctx context.Context, businessObjects ...string) (string, error) {
	return s.cntrInterface.Next(ctx, "", businessObjects)
}
