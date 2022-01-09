package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"testing"
	"time"

	"github.com/google/uuid"
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

			v := vouchers[uuid.NewSHA1(uuid.Nil, []byte("0000"))]

			assertions.Len(vouchers, 1)
			assertions.Len(v.LineItems(), 2)
			assertions.Equal(d100, v.Credit())
			assertions.Equal(d100, v.Debit())
			assertions.Equal("0000", v.Creator())
			assertions.True(accServiceMock.invoked)
		})
	}
}

func createUpdateVoucherCmd() *UpdateVoucherCmd {
	lineItems := []LineItemCmd{
		{
			Id:            uuid.New(),
			Summary:       "test_item1",
			AccountNumber: "1000",
			Debit:         "100",
			Credit:        "",
		},
		{
			Id:            uuid.New(),
			Summary:       "test_item2",
			AccountNumber: "1000",
			Debit:         "",
			Credit:        "100",
		},
	}
	return &UpdateVoucherCmd{
		VoucherUUID: uuid.NewSHA1(uuid.Nil, []byte("0000")),
		LineItems:   lineItems,
	}
}

func (r voucherRepoMock) initTestData() {
	item0, _ := domain.NewLineItem(uuid.New(), uuid.New(), "test_item0", "10", "")
	item1, _ := domain.NewLineItem(uuid.New(), uuid.New(), "test_item1", "", "10")
	items := []*domain.LineItem{item0, item1}
	v, _ := domain.NewVoucher(
		uuid.NewSHA1(uuid.Nil, []byte("0000")),
		uuid.New(),
		"GENERAL_VOUCHER",
		"1",
		0,
		items,
		"0000",
		"",
		"",
		false,
		false,
		false,
		time.Now(),
	)
	r.vouchers[v.Id()] = *v
}
