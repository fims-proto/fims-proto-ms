package adapter

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
	"strconv"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter_MemoryRepository_Next(t *testing.T) {
	t.Parallel()
	counterUUID := uuid.New()
	repo := prepareMemoryRepo(counterUUID)

	target := struct {
		data map[string]bool
		lock *sync.RWMutex
	}{
		data: make(map[string]bool, 30),
		lock: &sync.RWMutex{},
	}

	for i := 1; i <= 30; i++ {
		target.data["记"+strconv.Itoa(i)] = false
	}

	scheduleRaceTest(30, func(i int) {
		target.lock.Lock()
		defer target.lock.Unlock()

		ident, err := repo.UpdateAndRead(
			context.Background(),
			counterUUID,
			func(c *counter.Counter) (*counter.Counter, interface{}, error) {
				c.Next()
				return c, c.Identifier(), nil
			},
		)
		require.NoError(t, err)
		require.IsType(t, "string", ident)
		assert.False(t, target.data[ident.(string)])
		target.data[ident.(string)] = true
	})
	for _, v := range target.data {
		assert.True(t, v)
	}
}

func prepareMemoryRepo(counterUUID uuid.UUID) *CounterMemoryRepository {
	repo := NewCounterMemoryRepository()
	counter, _ := counter.NewCounter(counterUUID, "DUMMY", "记", "")
	_ = repo.CreateCounter(context.Background(), counter)
	return &repo
}

func scheduleRaceTest(workers int, startToDo func(i int)) {
	workersDone := sync.WaitGroup{}
	workersDone.Add(workers)

	startWorkers := make(chan struct{})

	for i := 0; i < workers; i++ {
		i := i
		go func() {
			defer workersDone.Done()
			<-startWorkers
			startToDo(i)
		}()
	}

	close(startWorkers)
	workersDone.Wait()
}

func TestAdapter_MemoryRepository_ReadByBusinessObject(t *testing.T) {
	t.Parallel()

	counterUUID1 := uuid.New()
	counterUUID2 := uuid.New()
	counter1, _ := counter.NewCounter(counterUUID1, "BO_1", "BO_1", "")
	counter2, _ := counter.NewCounter(counterUUID2, "BO_2", "BO_2", "")

	repo := NewCounterMemoryRepository()
	_ = repo.CreateCounter(context.Background(), counter1)
	_ = repo.CreateCounter(context.Background(), counter2)

	counterQuery1, err := repo.CounterByBusinessObject(context.Background(), "BO_1")
	require.NoError(t, err)
	assert.Equal(t, counterUUID1, counterQuery1.CounterUUID)

	counterQuery2, err := repo.CounterByBusinessObject(context.Background(), "BO_2")
	require.NoError(t, err)
	assert.Equal(t, counterUUID2, counterQuery2.CounterUUID)

	_, err = repo.CounterByBusinessObject(context.Background(), "NOT_EXIST")
	assert.Error(t, err)
}
