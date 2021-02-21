package adapter

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter_MemoryRepository_Next(t *testing.T) {
	t.Parallel()
	repo := prepareMemoryRepo(10)

	scheduleRaceTest(10, func(i int) {
		next, err := repo.GetNextFromCounter(context.Background(), "test"+strconv.FormatInt(int64(i), 10))
		require.NoError(t, err)
		assert.Equal(t, "haha000001hehe", next)
	})
}

func prepareMemoryRepo(num int) CounterMemoryRepository {
	repo := NewCounterMemoryRepository()
	for i := 0; i < num; i++ {
		// fmt.Println(i)
		counter, _ := counter.NewCounter("test"+strconv.FormatInt(int64(i), 10), 6, "haha", "hehe")
		repo.AddCounter(context.Background(), counter)
	}
	return repo
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
