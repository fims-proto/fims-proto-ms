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

func TestAdapter_MemoryRepository_ReadAll(t *testing.T) {
	t.Parallel()
	repo := NewVoucherMemoryRepository()
	v, err := voucher.NewVoucher("0000", 1, time.Now(), 0, []lineitem.LineItem{}, "0000")
	require.NoError(t, err)
	repo.data["0000"] = *v
	for i := 0; i < 100; i++ {
		go func() {
			vouchers, err := repo.AllVouchers(context.Background())
			require.NoError(t, err)
			assert.Equal(t, 1, len(vouchers))
		}()
	}
}

func TestAdapter_MemoryRepository_Add(t *testing.T) {
	t.Parallel()
	repo := NewVoucherMemoryRepository()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := voucher.NewVoucher(strconv.FormatInt(int64(i), 10), 1, time.Now(), 0, []lineitem.LineItem{}, "0000")
			require.NoError(t, err)
			err = repo.AddVoucher(context.Background(), v)
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	assert.Equal(t, 100, len(repo.data))
}

func TestAdapter_MemoryRepository_Update(t *testing.T) {
	t.Parallel()
	repo := NewVoucherMemoryRepository()
	v, err := voucher.NewVoucher("0000", 1, time.Now(), 0, []lineitem.LineItem{}, "0000")
	require.NoError(t, err)
	repo.data["0000"] = *v
	for i := 0; i < 100; i++ {
		go func() {
			err := repo.UpdateVoucher(context.Background(), "0000", func(v *voucher.Voucher) (*voucher.Voucher, error) {
				return v, nil// do nothing
			})
			require.NoError(t, err)
		}()
	}
}
