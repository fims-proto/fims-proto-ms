package counter

import (
	"context"
	counterport "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraprocessAdapter struct {
	cntrInterface counterport.CounterInterface
}

func NewIntraprocessAdapter(cntrInterface counterport.CounterInterface) IntraprocessAdapter {
	return IntraprocessAdapter{cntrInterface: cntrInterface}
}

func (s IntraprocessAdapter) Next(ctx context.Context, counterUUID uuid.UUID) (string, error) {
	return s.cntrInterface.Next(ctx, counterUUID)
}
