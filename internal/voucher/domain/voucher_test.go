package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sobId = uuid.New()

func TestDomain_NewVoucher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		voucherType string
		id          uuid.UUID
		sobId       uuid.UUID
		number      string
		items       []*LineItem
		verify      func(t *testing.T, voucher *Voucher, err error)
	}{
		{
			"normal_success", "GENERAL_VOUCHER", uuid.NewSHA1(uuid.Nil, []byte("test_uuid")), sobId, "1", prepareBalancedItems(),
			func(t *testing.T, voucher *Voucher, err error) {
				require.NoError(t, err)
				assert.Equal(t, uuid.NewSHA1(uuid.Nil, []byte("test_uuid")), voucher.Id())
				assert.Equal(t, "1", voucher.Number())
				assert.Equal(t, "200", voucher.Credit().String())
				assert.Equal(t, "200", voucher.Debit().String())
			},
		},
		{
			"imbalanced_error", "GENERAL_VOUCHER", uuid.New(), sobId, "1", prepareImbalancedItems(),
			func(t *testing.T, voucher *Voucher, err error) {
				require.Nil(t, voucher)
				assert.Error(t, err)
			},
		},
		{
			"empty_line_item_error", "GENERAL_VOUCHER", uuid.New(), sobId, "1",
			[]*LineItem{},
			func(t *testing.T, voucher *Voucher, err error) {
				require.Nil(t, voucher)
				assert.Error(t, err)
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			voucher, err := NewVoucher(test.id, test.sobId, test.voucherType, test.number, 0, test.items, "creator", "", "", false, false, false, time.Now())
			test.verify(t, voucher, err)
		})
	}
}

func prepareBalancedItems() []*LineItem {
	accountId := uuid.New()
	item1, _ := NewLineItem(uuid.New(), accountId, "1000", decimal.RequireFromString("100"), decimal.Zero)
	item2, _ := NewLineItem(uuid.New(), accountId, "1001", decimal.RequireFromString("100"), decimal.Zero)
	item3, _ := NewLineItem(uuid.New(), accountId, "2000", decimal.Zero, decimal.RequireFromString("150"))
	item4, _ := NewLineItem(uuid.New(), accountId, "2001", decimal.Zero, decimal.RequireFromString("50"))
	return []*LineItem{
		item1,
		item2,
		item3,
		item4,
	}
}

func prepareImbalancedItems() []*LineItem {
	accountId := uuid.New()
	item1, _ := NewLineItem(uuid.New(), accountId, "1000", decimal.RequireFromString("100"), decimal.Zero)
	item2, _ := NewLineItem(uuid.New(), accountId, "1001", decimal.RequireFromString("200"), decimal.Zero)
	item3, _ := NewLineItem(uuid.New(), accountId, "2000", decimal.Zero, decimal.RequireFromString("150"))
	item4, _ := NewLineItem(uuid.New(), accountId, "2001", decimal.Zero, decimal.RequireFromString("50"))
	return []*LineItem{
		item1,
		item2,
		item3,
		item4,
	}
}
