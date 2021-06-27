package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_NewVoucher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		voucherType VoucherType
		uuid        uuid.UUID
		number      string
		items       []LineItem
		verify      func(t *testing.T, voucher *Voucher, err error)
	}{
		{
			"normal_success", GeneralVoucher, uuid.NewSHA1(uuid.Nil, []byte("test_uuid")), "1", prepareBalancedItems(),
			func(t *testing.T, voucher *Voucher, err error) {
				require.NoError(t, err)
				assert.Equal(t, uuid.NewSHA1(uuid.Nil, []byte("test_uuid")), voucher.UUID())
				assert.Equal(t, "1", voucher.Number())
				assert.Equal(t, "200", voucher.Credit().String())
				assert.Equal(t, "200", voucher.Debit().String())
			},
		},
		{
			"imbalanced_error", GeneralVoucher, uuid.New(), "1", prepareImbalancedItems(),
			func(t *testing.T, voucher *Voucher, err error) {
				require.Nil(t, voucher)
				assert.Error(t, err)
			},
		},
		{
			"empty_lineitem_error", GeneralVoucher, uuid.New(), "1",
			[]LineItem{},
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
			voucher, err := NewVoucher("test_sob", test.uuid, test.voucherType, test.number, time.Now(), 0, test.items, "")
			test.verify(t, voucher, err)
		})
	}
}

func prepareBalancedItems() []LineItem {
	item1, _ := NewLineItem("test", "1000", "100", "")
	item2, _ := NewLineItem("test", "1001", "100", "")
	item3, _ := NewLineItem("test", "2000", "", "150")
	item4, _ := NewLineItem("test", "2001", "", "50")
	return []LineItem{
		*item1,
		*item2,
		*item3,
		*item4,
	}
}

func prepareImbalancedItems() []LineItem {
	item1, _ := NewLineItem("test", "1000", "100", "")
	item2, _ := NewLineItem("test", "1001", "200", "")
	item3, _ := NewLineItem("test", "2000", "", "150")
	item4, _ := NewLineItem("test", "2001", "", "50")
	return []LineItem{
		*item1,
		*item2,
		*item3,
		*item4,
	}
}
