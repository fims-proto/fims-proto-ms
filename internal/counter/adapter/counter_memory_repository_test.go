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
	counter, _ := counter.NewCounter(counterUUID, "", "记", "")
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
