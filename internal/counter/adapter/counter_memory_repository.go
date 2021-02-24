package adapter

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CounterWrapper struct {
	Counter *counter.Counter
	lock    *sync.RWMutex
}

type CounterMemoryRepository struct {
	data sync.Map
}

func NewCounterMemoryRepository() CounterMemoryRepository {
	return CounterMemoryRepository{
		data: sync.Map{},
	}
}

func (r *CounterMemoryRepository) CreateCounter(ctx context.Context, c *counter.Counter) error {
	_, ok := r.data.Load(c.UUID())
	if ok {
		return errors.Errorf("counter with UUID %s already exists", c.UUID())
	}
	r.data.Store(
		c.UUID(),
		&CounterWrapper{
			lock:    &sync.RWMutex{},
			Counter: c,
		},
	)
	return nil
}

func (r *CounterMemoryRepository) UpdateCounter(ctx context.Context, counterUUID uuid.UUID, updateFn func(c *counter.Counter) (*counter.Counter, error)) error {
	counterW, ok := r.data.Load(counterUUID)
	if !ok {
		return errors.Errorf("counter %s does not exist", counterUUID)
	}

	counterW.(*CounterWrapper).lock.Lock()
	defer counterW.(*CounterWrapper).lock.Unlock()

	c, err := updateFn(counterW.(*CounterWrapper).Counter)
	if err != nil {
		return errors.Wrapf(err, "counter %s update failed", counterUUID)
	}
	counterW.(*CounterWrapper).Counter = c
	r.data.Store(counterUUID, counterW)
	return nil
}

func (r *CounterMemoryRepository) UpdateAndRead(
	ctx context.Context,
	counterUUID uuid.UUID,
	updateAndReadFn func(c *counter.Counter) (*counter.Counter, interface{}, error),
) (interface{}, error) {
	counterW, ok := r.data.Load(counterUUID)
	if !ok {
		return nil, errors.Errorf("counter %s does not exist", counterUUID)
	}

	counterW.(*CounterWrapper).lock.Lock()
	defer counterW.(*CounterWrapper).lock.Unlock()

	c, readValue, err := updateAndReadFn(counterW.(*CounterWrapper).Counter)
	if err != nil {
		return nil, errors.Wrapf(err, "counter %s update and read failed", counterUUID)
	}
	counterW.(*CounterWrapper).Counter = c
	r.data.Store(counterUUID, counterW)
	return readValue, nil
}

func (r *CounterMemoryRepository) DeleteCounter(ctx context.Context, UUID string) error {
	r.data.Delete(UUID)
	return nil
}
