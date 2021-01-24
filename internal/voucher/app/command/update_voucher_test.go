package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestApp_HandleUpdateVoucherHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		constructor func() *UpdateVoucherCmd
	}{
		{
			name: "normal_success",
			constructor: func() *UpdateVoucherCmd {
				return createUpdateVoucherCmd()
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assertions := assert.New(t)

			cmd := test.constructor()
			repoMock := newVoucherRepoMock()
			repoMock.initTestData()
			accServiceMock := newAccountService()
			handler := NewUpdateVoucherHandler(repoMock, &accServiceMock)
			err := handler.Handle(context.Background(), *cmd)

			assertions.NoError(err)
			vouchers := repoMock.vouchers

			d100, _ := decimal.NewFromString("100")

			assertions.Len(vouchers, 1)
			assertions.Len(vouchers["0000"].LineItems(), 2)
			assertions.Equal(d100, vouchers["0000"].Credit())
			assertions.Equal(d100, vouchers["0000"].Debit())
			assertions.Equal("0000", vouchers["0000"].CreatorUUID())
			assertions.True(accServiceMock.invoked)
		})
	}
}

func createUpdateVoucherCmd() *UpdateVoucherCmd {
	lineItems := []LineItemCmd{
		{
			Summary:       "test_item1",
			AccountNumber: "1000",
			Debit:         "100",
			Credit:        "",
		},
		{
			Summary:       "test_item2",
			AccountNumber: "1000",
			Debit:         "",
			Credit:        "100",
		},
	}
	return &UpdateVoucherCmd{
		VoucherUUID: "0000",
		LineItems:   lineItems,
	}
}

func (r voucherRepoMock) initTestData() {
	item0, _ := lineitem.NewLineItem("test_item0", "1000", "10", "")
	item1, _ := lineitem.NewLineItem("test_item1", "1000", "", "10")
	items := []lineitem.LineItem{*item0, *item1}
	v, _ := voucher.NewVoucher(
		"0000", 1,
		time.Now(),
		0,
		items,
		"0000",
	)
	r.vouchers["0000"] = *v
}
