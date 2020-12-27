package voucher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"testing"
	"time"
)

func TestDomain_NewVoucher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		uuid   string
		number uint
		items  []lineitem.LineItem
		verify func(t *testing.T, voucher *Voucher, err error)
	}{
		{
			"normal_success", "test_uuid", 1, prepareBalancedItems(),
			func(t *testing.T, voucher *Voucher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "test_uuid", voucher.UUID())
				assert.Equal(t, uint(1), voucher.Number())
				assert.Equal(t, "200", voucher.Credit().String())
				assert.Equal(t, "200", voucher.Debit().String())
			},
		},
		{
			"imbalanced_error", "test_uuid", 1, prepareImbalancedItems(),
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
			voucher, err := NewVoucher(test.uuid, test.number, time.Now(), 0, test.items, "", "", false, "", false)
			test.verify(t, voucher, err)
		})
	}
}

func prepareBalancedItems() []lineitem.LineItem {
	item1, _ := lineitem.NewLineItem("test", "1000", "100", "")
	item2, _ := lineitem.NewLineItem("test", "1001", "100", "")
	item3, _ := lineitem.NewLineItem("test", "2000", "", "150")
	item4, _ := lineitem.NewLineItem("test", "2001", "", "50")
	return []lineitem.LineItem{
		*item1,
		*item2,
		*item3,
		*item4,
	}
}

func prepareImbalancedItems() []lineitem.LineItem {
	item1, _ := lineitem.NewLineItem("test", "1000", "100", "")
	item2, _ := lineitem.NewLineItem("test", "1001", "200", "")
	item3, _ := lineitem.NewLineItem("test", "2000", "", "150")
	item4, _ := lineitem.NewLineItem("test", "2001", "", "50")
	return []lineitem.LineItem{
		*item1,
		*item2,
		*item3,
		*item4,
	}
}
