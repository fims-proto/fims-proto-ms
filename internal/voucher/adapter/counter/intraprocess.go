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

func (s IntraprocessAdapter) Create(ctx context.Context, prefix string, sufix string) (uuid.UUID, error) {
	return s.cntrInterface.Create(ctx, prefix, sufix)
}

func (s IntraprocessAdapter) Next(ctx context.Context, counterUUID uuid.UUID) (string, error) {
	return s.cntrInterface.Next(ctx, counterUUID)
}

func (s IntraprocessAdapter) Delete(ctx context.Context, counterUUID uuid.UUID) error {
	return s.cntrInterface.Delete(ctx, counterUUID)
}

func (s IntraprocessAdapter) Reset(ctx context.Context, counterUUID uuid.UUID) error {
	return s.cntrInterface.Reset(ctx, counterUUID)
}
