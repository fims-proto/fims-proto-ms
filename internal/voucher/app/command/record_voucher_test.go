package command

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"testing"
	"time"
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
			deps := newRecordMockDeps()
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

type recordMockDeps struct {
	repository *recordRepoMock
	handler    RecordVoucherHandler
}

func newRecordMockDeps() recordMockDeps {
	repository := &recordRepoMock{vouchers: make(map[string]voucher.Voucher)}
	return recordMockDeps{
		repository: repository,
		handler:    RecordVoucherHandler{repository},
	}
}

type recordRepoMock struct {
	vouchers map[string]voucher.Voucher
}

func (r recordRepoMock) AddVoucher(ctx context.Context, v *voucher.Voucher) error {
	r.vouchers[v.UUID()] = *v
	return nil
}

func (r recordRepoMock) UpdateVoucher(ctx context.Context, voucherUUID string, updateFn func(v *voucher.Voucher) (*voucher.Voucher, error)) error {
	panic("implement me")
}
