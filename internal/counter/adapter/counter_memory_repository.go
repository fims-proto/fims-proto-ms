package adapter

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
	"sync"

	"github.com/pkg/errors"
)

type CounterMemoryRepository struct {
	lock *sync.RWMutex
	data map[string]counter.Counter
}

func NewCounterMemoryRepository() CounterMemoryRepository{
	return CounterMemoryRepository{
		lock: &sync.RWMutex{},
		data: make(map[string]counter.Counter),
	}
}

func (r *CounterMemoryRepository) ResetCounter(ctx context.Context, UUID string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	counter, ok := r.data[UUID]
	if !ok {
		errors.Errorf("Counter %s does not exist", UUID)
	}
	err := counter.Reset()
	if err != nil {
		return errors.Wrapf(err, "Counter %s reset failed", UUID)
	}
	r.data[UUID] = counter 
	return nil
}

func (r *CounterMemoryRepository) GetNextFromCounter(ctx context.Context, UUID string) (string,error){
	r.lock.Lock()
	defer r.lock.Unlock()
	counter, ok := r.data[UUID]
	if !ok {
		errors.Errorf("Counter %s does not exist", UUID)
	}
	next, err := counter.Next()
	if err != nil {
		return "",errors.Wrapf(err, "Counter %s next failed", UUID)
	}
	r.data[UUID] = counter
	return next, nil
}