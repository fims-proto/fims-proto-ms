package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"testing"
	"time"

	"github.com/pkg/errors"
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
			deps := newupdateDepsMock()
			err := deps.handler.Handle(context.Background(), *cmd)

			assertions.NoError(err)
			vouchers := deps.repository.vouchers

			d100, _ := decimal.NewFromString("100")

			assertions.Equal(1, len(vouchers))
			assertions.Equal(2, len(vouchers["0000"].LineItems()))
			assertions.Equal(d100, vouchers["0000"].Credit())
			assertions.Equal(d100, vouchers["0000"].Debit())
			assertions.Equal("0000", vouchers["0000"].CreatorUUID())
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
			Summary: "test_item2",
			AccountNumber: "1000",
			Debit: "",
			Credit: "100",
		},
	}
	return &UpdateVoucherCmd{
		UUID: "0000",
		LineItems:     lineItems,
	}
}

type updateDepsMock struct {
	repository *updateRepoMock
	handler    UpdateVoucherHandler
}

func newupdateDepsMock() updateDepsMock {
	repository := &updateRepoMock{vouchers: make(map[string]voucher.Voucher)}
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
	repository.AddVoucher(context.Background(), v)
	return updateDepsMock{
		repository: repository,
		handler:    UpdateVoucherHandler{repository},
	}
}

type updateRepoMock struct {
	vouchers map[string]voucher.Voucher
}

func (h *updateRepoMock) AddVoucher(ctx context.Context, voucher *voucher.Voucher) error {

	_, ok := h.vouchers[voucher.UUID()]
	if ok {
		return errors.Errorf("voucher %s exists", voucher.UUID())
	}

	h.vouchers[voucher.UUID()] = *voucher
	return nil
}

func (h *updateRepoMock) UpdateVoucher(ctx context.Context, voucherUUID string, updateFn func(v *voucher.Voucher) (*voucher.Voucher, error)) error {

	v, ok := h.vouchers[voucherUUID]
	if !ok {
		return errors.Errorf("voucher %s not exists", voucherUUID)
	}

	updatedVoucher, err := updateFn(&v)
	if err != nil {
		return errors.Wrapf(err, "voucher %s updated failed", voucherUUID)
	}
	h.vouchers[voucherUUID] = *updatedVoucher
	return nil
}
