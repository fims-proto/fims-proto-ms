package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"
	"sync"

	"github.com/pkg/errors"
)

type SobMemoryRepository struct {
	lock *sync.RWMutex
	data map[string]domain.Sob
}

func NewSobMemoryRepository() SobMemoryRepository {
	return SobMemoryRepository{
		lock: &sync.RWMutex{},
		data: make(map[string]domain.Sob),
	}
}

func (r SobMemoryRepository) CreateSob(ctx context.Context, sob *domain.Sob) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.data[sob.Id()]; ok {
		return errors.Errorf("sob exists with id %s", sob.Id())
	}

	r.data[sob.Id()] = *sob
	return nil
}

func (r SobMemoryRepository) UpdateSob(
	ctx context.Context,
	sobId string,
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

func (r SobMemoryRepository) SobById(ctx context.Context, sobId string) (query.Sob, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	sob, ok := r.data[sobId]
	if !ok {
		return query.Sob{}, errors.Errorf("sob %s does not exist", sobId)
	}

	return mapFromDomainSob(sob), nil
}

func mapFromDomainSob(sob domain.Sob) query.Sob {
	return query.Sob{
		Id:          sob.Id(),
		Name:        sob.Name(),
		Description: sob.Description(),
	}
}
