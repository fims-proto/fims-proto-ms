package counter

import (
	"context"
	counterport "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"
)

type IntraprocessService struct {
	cntrInterface counterport.CounterInterface
}

func NewIntraprocessService(cntrInterface counterport.CounterInterface) IntraprocessService {
	return IntraprocessService{cntrInterface: cntrInterface}
}

func (s IntraprocessService) Add(ctx context.Context, UUID string, len uint, prefix string, sufix string) error {
	return s.cntrInterface.Add(ctx, UUID, len, prefix, sufix)
}

func (s IntraprocessService) Next(ctx context.Context, UUID string) (string, error) {
	return s.cntrInterface.Next(ctx, UUID)
}

func (s IntraprocessService) Delete(ctx context.Context, UUID string) error {
	return s.cntrInterface.Delete(ctx, UUID)
}

func (s IntraprocessService) Reset(ctx context.Context, UUID string) error {
	return s.cntrInterface.Reset(ctx, UUID)
}
