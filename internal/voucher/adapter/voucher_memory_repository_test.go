package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"strconv"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// this is a simple check that the adapter implements the domain interface
func TestAdapter_MemoryRepository_InterfaceImplemented(t *testing.T) {
	t.Parallel()
	var _ domain.Repository = (*VoucherMemoryRepository)(nil)
}

func TestAdapter_MemoryRepository_ReadAll(t *testing.T) {
	t.Parallel()
	voucherUUID := uuid.New()
	repo := prepareMemoryRepo(t, "test_sob", voucherUUID)

	scheduleRaceTest(30, func(_ int) {
		vouchers, err := repo.AllVouchers(context.Background(), "test_sob")
		require.NoError(t, err)
		assert.Len(t, vouchers, 1)
	})
}

func TestAdapter_MemoryRepository_ReadOne(t *testing.T) {
	t.Parallel()
	voucherUUID := uuid.New()
	repo := prepareMemoryRepo(t, "test_sob", voucherUUID)

	scheduleRaceTest(30, func(_ int) {
		v, err := repo.VoucherByUUID(context.Background(), "test_sob", voucherUUID)
		require.NoError(t, err)
		assert.Equal(t, voucherUUID, v.UUID)
	})
}

func TestAdapter_MemoryRepository_Add(t *testing.T) {
	t.Parallel()
	repo := NewVoucherMemoryRepository()

	scheduleRaceTest(30, func(i int) {
		v, err := domain.NewVoucher("test_sob", uuid.New(), domain.GeneralVoucher, strconv.Itoa(i), 0, prepareBalancedItems(), "0000")
		require.NoError(t, err)
		_, err = repo.AddVoucher(context.Background(), v)
		require.NoError(t, err)
	})

	assert.Len(t, repo.data, 30)
}

func TestAdapter_MemoryRepository_Update(t *testing.T) {
	t.Parallel()
	voucherUUID := uuid.New()
	repo := prepareMemoryRepo(t, "test_sob", voucherUUID)

	voucherAudited := make(chan int, 30)

	scheduleRaceTest(30, func(i int) {
		err := repo.UpdateVoucher(context.Background(), "test_sob", voucherUUID, func(v *domain.Voucher) (*domain.Voucher, error) {
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

func prepareMemoryRepo(t *testing.T, sob string, voucherUUID uuid.UUID) VoucherMemoryRepository {
	repo := NewVoucherMemoryRepository()
	v, err := domain.NewVoucher(sob, voucherUUID, domain.GeneralVoucher, "1", 0, prepareBalancedItems(), "0000")
	require.NoError(t, err)
	repo.data[sob+v.UUID().String()] = *v
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

func prepareBalancedItems() []domain.LineItem {
	item1, _ := domain.NewLineItem(uuid.New(), "test", "1000", "100", "")
	item2, _ := domain.NewLineItem(uuid.New(), "test", "1001", "100", "")
	item3, _ := domain.NewLineItem(uuid.New(), "test", "2000", "", "150")
	item4, _ := domain.NewLineItem(uuid.New(), "test", "2001", "", "50")
	return []domain.LineItem{
		*item1,
		*item2,
		*item3,
		*item4,
	}
}
