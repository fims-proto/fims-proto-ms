package adapter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"strconv"
	"sync"
	"testing"
	"time"
)

// this is a simple check that the adapter implements the domain interface
func TestAdapter_MemoryRepository_InterfaceImplemented(t *testing.T){
	t.Parallel()
	var _ voucher.Repository = (*VoucherMemoryRepository)(nil)
}

func TestAdapter_MemoryRepository_ReadAll(t *testing.T) {
	t.Parallel()
	repo := prepareMemoryRepo(t)

	scheduleRaceTest(20, func(_ int) {
		vouchers, err := repo.AllVouchers(context.Background())
		require.NoError(t, err)
		assert.Len(t, vouchers, 1)
	})
}

func TestAdapter_MemoryRepository_ReadOne(t *testing.T) {
	t.Parallel()
	repo := prepareMemoryRepo(t)

	scheduleRaceTest(20, func(_ int) {
		v, err := repo.VoucherForUUID("0000", context.Background())
		require.NoError(t, err)
		assert.Equal(t, "0000", v.UUID)
	})
}

func TestAdapter_MemoryRepository_Add(t *testing.T) {
	t.Parallel()
	repo := NewVoucherMemoryRepository()

	scheduleRaceTest(20, func(i int) {
		v, err := voucher.NewVoucher(strconv.FormatInt(int64(i), 10), 1, time.Now(), 0, []lineitem.LineItem{}, "0000")
		require.NoError(t, err)
		err = repo.AddVoucher(context.Background(), v)
		require.NoError(t, err)
	})

	assert.Len(t, repo.data, 20)
}

func TestAdapter_MemoryRepository_Update(t *testing.T) {
	t.Parallel()
	repo := prepareMemoryRepo(t)

	voucherAudited := make(chan int, 20)

	scheduleRaceTest(20, func(i int) {
		err := repo.UpdateVoucher(context.Background(), "0000", func(v *voucher.Voucher) (*voucher.Voucher, error) {
			if err := v.Audit("testUUID"); err == nil {
				// success
				voucherAudited <- i
			}
			return v, nil
		})
		require.NoError(t, err)
	})

	assert.Len(t, voucherAudited, 1, "voucher should be audit only once")
}

func prepareMemoryRepo(t *testing.T) VoucherMemoryRepository {
	repo := NewVoucherMemoryRepository()
	v, err := voucher.NewVoucher("0000", 1, time.Now(), 0, []lineitem.LineItem{}, "0000")
	require.NoError(t, err)
	repo.data["0000"] = *v
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
