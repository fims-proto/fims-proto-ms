package counter

import (
	"context"
	counterport "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraprocessService struct {
	cntrInterface counterport.CounterInterface
}

func NewIntraprocessService(cntrInterface counterport.CounterInterface) IntraprocessService {
	return IntraprocessService{cntrInterface: cntrInterface}
}

func (s IntraprocessService) Create(ctx context.Context, prefix string, sufix string) (uuid.UUID, error) {
	return s.cntrInterface.Create(ctx, prefix, sufix)
}

func (s IntraprocessService) Next(ctx context.Context, counterUUID uuid.UUID) (string, error) {
	return s.cntrInterface.Next(ctx, counterUUID)
}

func (s IntraprocessService) Delete(ctx context.Context, counterUUID uuid.UUID) error {
	return s.cntrInterface.Delete(ctx, counterUUID)
}

func (s IntraprocessService) Reset(ctx context.Context, counterUUID uuid.UUID) error {
	return s.cntrInterface.Reset(ctx, counterUUID)
}
