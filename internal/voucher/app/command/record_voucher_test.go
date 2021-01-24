package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestApp_HandleRecordVoucherHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		constructor func() *RecordVoucherCmd
	}{
		{
			name: "normal_success",
			constructor: func() *RecordVoucherCmd {
				return createVoucherCmd()
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
			accServiceMock := newAccountService()
			handler := NewRecordVoucherHandler(repoMock, &accServiceMock)
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

func createVoucherCmd() *RecordVoucherCmd {
	lineItems := []LineItemCmd{
		{
			Summary:       "test_item1",
			AccountNumber: "1000",
			Debit:         "100",
			Credit:        "",
		},
		{
			Summary:       "test_item2",
			AccountNumber: "2000",
			Debit:         "",
			Credit:        "100",
		},
	}
	return &RecordVoucherCmd{
		UUID:               "0000",
		Number:             1,
		CreatedAt:          time.Now(),
		AttachmentQuantity: 0,
		LineItems:          lineItems,
		Debit:              "100",
		Credit:             "100",
		CreatorUUID:        "0000",
	}
}

func newVoucherRepoMock() voucherRepoMock {
	return voucherRepoMock{vouchers: make(map[string]voucher.Voucher)}
}

func newAccountService() accountServiceMock {
	return accountServiceMock{invoked: false}
}

type voucherRepoMock struct {
	vouchers map[string]voucher.Voucher
}

func (r voucherRepoMock) AddVoucher(ctx context.Context, v *voucher.Voucher) error {
	r.vouchers[v.UUID()] = *v
	return nil
}

func (r voucherRepoMock) UpdateVoucher(ctx context.Context, voucherUUID string, updateFn func(v *voucher.Voucher) (*voucher.Voucher, error)) error {
	v, ok := r.vouchers[voucherUUID]
	if !ok {
		return errors.Errorf("voucher %s not exists", voucherUUID)
	}

	updatedVoucher, err := updateFn(&v)
	if err != nil {
		return errors.Wrapf(err, "voucher %s updated failed", voucherUUID)
	}
	r.vouchers[voucherUUID] = *updatedVoucher
	return nil
}

type accountServiceMock struct {
	invoked bool
}

func (s *accountServiceMock) ValidateExistence(ctx context.Context, accNumbers []string) error {
	s.invoked = true
	return nil
}
