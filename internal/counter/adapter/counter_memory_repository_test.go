package adapter

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
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

	target := map[string]bool{
		"记1":  false,
		"记2":  false,
		"记3":  false,
		"记4":  false,
		"记5":  false,
		"记6":  false,
		"记7":  false,
		"记8":  false,
		"记9":  false,
		"记10": false,
	}

	scheduleRaceTest(10, func(i int) {
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
		assert.False(t, target[ident.(string)])
		target[ident.(string)] = true
	})
	for _, v := range target {
		assert.True(t, v)
	}
}

func prepareMemoryRepo(counterUUID uuid.UUID) *CounterMemoryRepository {
	repo := NewCounterMemoryRepository()
	counter, _ := counter.NewCounter(counterUUID, "记", "")
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
