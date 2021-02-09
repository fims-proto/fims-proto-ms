package adapter

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
	"sync"

	"github.com/pkg/errors"
)

type CounterWrapper struct {
	Counter *counter.Counter
	lock *sync.RWMutex
}

type CounterMemoryRepository struct {
	data sync.Map
}

func NewCounterMemoryRepository() CounterMemoryRepository{
	return CounterMemoryRepository{
		data: sync.Map{},
	}
}

func (r *CounterMemoryRepository) AddCounter(ctx context.Context, counter *counter.Counter) error {
	_, ok := r.data.Load(counter.UUID)
	if ok {
		return errors.Errorf("Counter with UUID %s already exists", counter.UUID)
	}
	r.data.Store(	
		counter.UUID, 
		CounterWrapper{
			lock: &sync.RWMutex{},
			Counter: counter,
		},
	)
	return nil
}

func (r *CounterMemoryRepository) ResetCounter(ctx context.Context, UUID string) error {
	counterW, ok := r.data.Load(UUID)
	if !ok {
		return errors.Errorf("Counter %s does not exist", UUID)
	}
	counterW.(CounterWrapper).lock.Lock()
	defer counterW.(CounterWrapper).lock.Unlock()
	err := counterW.(CounterWrapper).Counter.Reset()
	if err != nil {
		return errors.Wrapf(err, "Counter %s reset failed", UUID)
	}
	r.data.Store(UUID, counterW) 
	return nil
}

func (r *CounterMemoryRepository) GetNextFromCounter(ctx context.Context, UUID string) (string,error){
	counterW, ok := r.data.Load(UUID)
	if !ok {
		return "", errors.Errorf("Counter %s does not exist", UUID)
	}
	next, err := counterW.(CounterWrapper).Counter.Next()
	if err != nil {
		return "",errors.Wrapf(err, "Counter %s next failed", UUID)
	}
	r.data.Store(UUID,counterW)
	return next, nil
}