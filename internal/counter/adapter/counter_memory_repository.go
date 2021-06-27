package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/app/query"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"
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

func (r *CounterMemoryRepository) DeleteCounter(ctx context.Context, counterUUID uuid.UUID) error {
	r.data.Delete(counterUUID)
	return nil
}

func (r *CounterMemoryRepository) CounterByBusinessObject(ctx context.Context, sep string, businessObjects []string) (query.Counter, error) {
	m, err := domain.NewMatcher(sep, businessObjects...)
	if err != nil {
		return query.Counter{}, errors.Errorf("cannot create counter matcher: %s", businessObjects)
	}

	var counterUUID uuid.UUID
	r.data.Range(func(key, value interface{}) bool {
		if value.(*CounterWrapper).Counter.BusinessObject() == m.String() {
			counterUUID = value.(*CounterWrapper).Counter.UUID()
			return false
		}
		return true
	})
	if counterUUID == uuid.Nil {
		return query.Counter{}, errors.Errorf("cannot find counter with business object %s", m.String())
	}
	counterW, ok := r.data.Load(counterUUID)
	if !ok {
		return query.Counter{}, errors.Errorf("counter %s does not exist", counterUUID)
	}
	return MapFromDomainCounter(*counterW.(*CounterWrapper).Counter), nil
}

func MapFromDomainCounter(c domain.Counter) query.Counter {
	return query.Counter{
		CounterUUID: c.UUID(),
	}
}
