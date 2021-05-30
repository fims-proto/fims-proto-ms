package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type SobMemoryRepository struct {
	lock *sync.RWMutex
	data map[uuid.UUID]domain.Sob
}

func NewSobMemoryRepository() SobMemoryRepository {
	return SobMemoryRepository{
		lock: &sync.RWMutex{},
		data: make(map[uuid.UUID]domain.Sob),
	}
}

func (r SobMemoryRepository) CreateSob(ctx context.Context, sob *domain.Sob) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, dbSob := range r.data {
		if dbSob.Name() == sob.Name() {
			return errors.Errorf("sob exists with name %s", sob.Name())
		}
	}

	r.data[sob.UUID()] = *sob
	return nil
}

func (r SobMemoryRepository) UpdateSob(
	ctx context.Context,
	sobUUID uuid.UUID,
	updateFn func(s *domain.Sob) (*domain.Sob, error),
) {
	panic("not implemented")
}

func (r SobMemoryRepository) AllSobs(ctx context.Context) ([]query.Sob, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	sobs := []query.Sob{}
	for _, sob := range r.data {
		sobs = append(sobs, mapFromDomainSob(sob))
	}

	return sobs, nil
}

func mapFromDomainSob(sob domain.Sob) query.Sob {
	return query.Sob{
		UUID: sob.UUID(),
		Name: sob.Name(),
	}
}
